package stats

import (
	"context"
	"fmt"

	"github.com/karanbihani/file-vault/internal/db" // Adjust to your module path
)

// Service handles the business logic for statistics.
type Service struct {
	queries *db.Queries
}

// NewService creates a new stats service.
func NewService(queries *db.Queries) *Service {
	return &Service{
		queries: queries,
	}
}

// UserDashboardStatsResponse defines the structure of our comprehensive stats response.
type UserDashboardStatsResponse struct {
	FilesUploadedCount       int64   `json:"files_uploaded_count"`
	TotalDownloadsOnShares   int64   `json:"total_downloads_on_shares"`
	PublicSharesCount        int64   `json:"public_shares_count"`
	PrivateSharesCount       int64   `json:"private_shares_count"`
	DeduplicatedStorageUsage int64   `json:"deduplicated_storage_usage_bytes"`
	OriginalStorageUsage     int64   `json:"original_storage_usage_bytes"`
	StorageSavingsBytes      int64   `json:"storage_savings_bytes"`
	StorageSavingsPercentage float64 `json:"storage_savings_percentage"`
}

// GetUserDashboardStats calculates and returns the comprehensive statistics for a user.
func (s *Service) GetUserDashboardStats(ctx context.Context, userID int64) (*UserDashboardStatsResponse, error) {
	stats, err := s.queries.GetUserDashboardStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats from db: %w", err)
	}

	savings := stats.OriginalStorageUsage - stats.DeduplicatedStorageUsage
	var percentage float64
	if stats.OriginalStorageUsage > 0 {
		percentage = (float64(savings) / float64(stats.OriginalStorageUsage)) * 100
	}

	return &UserDashboardStatsResponse{
		FilesUploadedCount:       stats.FilesUploadedCount,
		TotalDownloadsOnShares:   stats.TotalDownloadsOnShares,
		PublicSharesCount:        stats.PublicSharesCount,
		PrivateSharesCount:       stats.PrivateSharesCount,
		DeduplicatedStorageUsage: stats.DeduplicatedStorageUsage,
		OriginalStorageUsage:     stats.OriginalStorageUsage,
		StorageSavingsBytes:      savings,
		StorageSavingsPercentage: percentage,
	}, nil
}