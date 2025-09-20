package admin

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

func (s *Service) ListAllFiles(ctx context.Context) ([]db.ListAllFilesRow, error) {
	return s.queries.ListAllFiles(ctx)
}

func (s *Service) GetSystemStats(ctx context.Context) (db.GetSystemStatsRow, error) {
	return s.queries.GetSystemStats(ctx)
}