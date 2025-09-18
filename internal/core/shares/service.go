package shares

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"

	"github.com/karanbihani/file-vault/internal/db"      // Adjust to your module path
	"github.com/karanbihani/file-vault/internal/storage" // Adjust to your module path
	"github.com/jackc/pgx/v5"
)

// Service handles the business logic for file sharing.
type Service struct {
	queries *db.Queries
	storage *storage.Client
}

// NewService creates a new sharing service.
func NewService(queries *db.Queries, storageClient *storage.Client) *Service {
	return &Service{
		queries: queries,
		storage: storageClient,
	}
}

// generateShareToken creates a cryptographically secure, random token.
func generateShareToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// CreatePublicLink creates a new public share link for a file.
func (s *Service) CreatePublicLink(ctx context.Context, fileID, ownerID int64) (*db.Share, error) {
	// Verify the user owns the file they are trying to share.
	_, err := s.queries.GetUserFileForDownload(ctx, db.GetUserFileForDownloadParams{ID: fileID, OwnerID: ownerID})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("file not found or access denied")
		}
		return nil, fmt.Errorf("failed to verify file ownership: %w", err)
	}

	token, err := generateShareToken(16)
	if err != nil {
		return nil, fmt.Errorf("failed to generate share token: %w", err)
	}

	share, err := s.queries.CreatePublicShareLink(ctx, db.CreatePublicShareLinkParams{
		UserFileID: fileID,
		ShareToken: token,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create share link in database: %w", err)
	}

	return &share, nil
}

// PublicDownloadResponse holds the data for a public file download.
type PublicDownloadResponse struct {
	Data     io.ReadCloser
	Filename string
	Size     int64
}

// ProcessPublicDownload verifies a token, gets the file, and increments the download count.
func (s *Service) ProcessPublicDownload(ctx context.Context, token string) (*PublicDownloadResponse, error) {
	shareMeta, err := s.queries.GetShareByToken(ctx, token)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("invalid or expired share link")
		}
		return nil, fmt.Errorf("failed to retrieve share link: %w", err)
	}

	object, err := s.storage.Get(ctx, shareMeta.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve file from storage: %w", err)
	}

	// Increment the download count in a background goroutine so it doesn't slow down the user.
	go func() {
		err := s.queries.IncrementShareDownloadCount(context.Background(), shareMeta.ID)
		if err != nil {
			log.Printf("ERROR: failed to increment download count for share ID %d: %v", shareMeta.ID, err)
		}
	}()

	return &PublicDownloadResponse{
		Data:     object,
		Filename: shareMeta.Filename,
		Size:     shareMeta.SizeBytes,
	}, nil
}