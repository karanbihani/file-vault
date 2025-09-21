package search

import (
	"context"
	"github.com/karanbihani/file-vault/internal/db"
)

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

// Pass through to the database query
func (s *Service) GetUserPermissions(ctx context.Context, userID int64) ([]string, error) {
	return s.queries.GetUserPermissions(ctx, userID)
}

// SearchFiles converts API parameters into the format required by the sqlc query.
func (s *Service) SearchFiles(ctx context.Context, params db.SearchFilesParams) ([]db.SearchFilesRow, error) {
	// Add wildcard '%' for ILIKE search on filename
	if params.Filename.Valid {
		params.Filename.String = "%" + params.Filename.String + "%"
	}
	return s.queries.SearchFiles(ctx, params)
}