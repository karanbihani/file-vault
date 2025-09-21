package api

import (
	"net/http"
	"strings"
	"time"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/karanbihani/file-vault/internal/core/search"
	"github.com/karanbihani/file-vault/internal/db"
	"github.com/karanbihani/file-vault/internal/auth" 
	"github.com/jackc/pgx/v5/pgtype"
)

type SearchHandler struct {
	searchService *search.Service
}

func NewSearchHandler(service *search.Service) *SearchHandler {
	return &SearchHandler{searchService: service}
}

func (h *SearchHandler) Search(c *gin.Context) {

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}

	// --- RBAC Logic ---
	permissions, err := h.searchService.GetUserPermissions(c.Request.Context(), userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve permissions"})
		return
	}

	isAdminSearch := false
	canSearch := false
	for _, p := range permissions {
		if p == auth.PermissionSearchAll {
			isAdminSearch = true
			canSearch = true
			break
		}
		if p == auth.PermissionSearchSelf {
			canSearch = true
		}
	}

	if !canSearch {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied: you do not have permission to search files"})
		return
	}
	// --- End RBAC Logic ---

	var params db.SearchFilesParams
	params.RequestingUserID = userID.(int64)
	params.IsAdmin = isAdminSearch
	
	// This handler carefully checks each query parameter. If it exists,
	// it's added to the params struct. Otherwise, the field remains nil
	// and the SQL query will ignore it.

	if filename, ok := c.GetQuery("filename"); ok {
		params.Filename = pgtype.Text{String: filename, Valid: true}
	}
	if mimeType, ok := c.GetQuery("mime_type"); ok {
		params.MimeType = pgtype.Text{String: mimeType, Valid: true}
	}
	if minSize, ok := c.GetQuery("min_size"); ok {
		if size, err := strconv.ParseInt(minSize, 10, 64); err == nil {
			params.MinSize = pgtype.Int8{Int64: size, Valid: true}
		}
	}
	if maxSize, ok := c.GetQuery("max_size"); ok {
		if size, err := strconv.ParseInt(maxSize, 10, 64); err == nil {
			params.MaxSize = pgtype.Int8{Int64: size, Valid: true}
		}
	}
	if startDate, ok := c.GetQuery("start_date"); ok {
		if t, err := time.Parse(time.RFC3339, startDate); err == nil {
			params.StartDate = pgtype.Timestamptz{Time: t, Valid: true}
		}
	}
	if endDate, ok := c.GetQuery("end_date"); ok {
		if t, err := time.Parse(time.RFC3339, endDate); err == nil {
			params.EndDate = pgtype.Timestamptz{Time: t, Valid: true}
		}
	}
	if tags, ok := c.GetQuery("tags"); ok {
		params.Tags = strings.Split(tags, ",")
	}
	if uploader, ok := c.GetQuery("uploader"); ok {
		params.UploaderEmail = pgtype.Text{String: uploader, Valid: true}
	}

	results, err := h.searchService.SearchFiles(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}