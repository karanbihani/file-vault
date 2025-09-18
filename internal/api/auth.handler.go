package api

import (
	"net/http"

	"github.com/karanbihani/file-vault/internal/auth" // Adjust to your module path
	"github.com/gin-gonic/gin"
)

// AuthHandler holds the dependencies for the auth handlers.
type AuthHandler struct {
	authService *auth.Service
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(service *auth.Service) *AuthHandler {
	return &AuthHandler{
		authService: service,
	}
}

// Register handles the POST /register endpoint.
func (h *AuthHandler) Register(c *gin.Context) {
	var params auth.RegisterUserParams
	// Bind the incoming JSON request body to our params struct.
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	user, err := h.authService.RegisterUser(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Don't send the password hash back to the client.
	user.PasswordHash = ""
	c.JSON(http.StatusCreated, user)
}

// Login handles the POST /login endpoint.
func (h *AuthHandler) Login(c *gin.Context) {
	var params auth.LoginUserParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	token, err := h.authService.LoginUser(c.Request.Context(), params)
	if err != nil {
		// Send a 401 Unauthorized status for invalid credentials.
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
