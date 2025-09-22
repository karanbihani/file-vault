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
	"github.com/karanbihani/file-vault/internal/core/audit"
	"github.com/jackc/pgx/v5"
)

// Service handles the business logic for file sharing.
type Service struct {
	queries *db.Queries
	storage *storage.Client
	auditService *audit.Service 
}

// NewService creates a new sharing service.
func NewService(queries *db.Queries, storageClient *storage.Client, auditService *audit.Service) *Service {
	return &Service{
		queries: queries,
		storage: storageClient,
		auditService: audit.NewService(queries),
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

// CreatePublicLink now includes a critical ownership verification check.
func (s *Service) CreatePublicLink(ctx context.Context, fileID, ownerID int64) (*db.Share, error) {
	// --- SECURITY FIX: Verify Ownership ---
	// Before creating a share link, we query the database to ensure the user making
	// the request is the actual owner of the file.
	// We can reuse the GetUserFileForDownload query as it performs this exact check.
	_, err := s.queries.GetUserFileForDownload(ctx, db.GetUserFileForDownloadParams{ID: fileID, OwnerID: ownerID})
	if err != nil {
		if err == pgx.ErrNoRows {
			// This error now correctly means "file not found OR you don't own it".
			return nil, fmt.Errorf("file not found or access denied")
		}
		// Handle other potential database errors.
		return nil, fmt.Errorf("failed to verify file ownership: %w", err)
	}
	// --- END SECURITY FIX ---

	// If the check above passes, we can safely proceed.
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

	s.auditService.LogActivity(ctx, ownerID, "share:create_public", map[string]interface{}{
		"file_id": fileID,
		"share_token": token,
	})
	
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

func (s *Service) ShareFileWithUser(ctx context.Context, fileID, ownerID int64, recipientEmail string) error {
	// 1. Verify the user owns the file they are trying to share.
	_, err := s.queries.GetUserFileForDownload(ctx, db.GetUserFileForDownloadParams{ID: fileID, OwnerID: ownerID})
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("file not found or access denied")
		}
		return fmt.Errorf("failed to verify file ownership: %w", err)
	}

	// 2. Find the recipient user by their email address.
	recipient, err := s.queries.GetUserByEmail(ctx, recipientEmail)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("recipient user with email '%s' not found", recipientEmail)
		}
		return fmt.Errorf("failed to find recipient user: %w", err)
	}

	// 3. Prevent users from sharing files with themselves.
	if ownerID == recipient.ID {
		return fmt.Errorf("cannot share a file with yourself")
	}

	// 4. Check if the file is already shared with this user to avoid duplicates.
	alreadyShared, err := s.queries.IsFileAlreadySharedWithUser(ctx, db.IsFileAlreadySharedWithUserParams{
		UserFileID:       fileID,
		SharedWithUserID: recipient.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to check for existing share: %w", err)
	}
	if alreadyShared {
		return fmt.Errorf("file is already shared with this user")
	}

	// 5. Create the share record in the database.
	err = s.queries.ShareFileWithUser(ctx, db.ShareFileWithUserParams{
		UserFileID:       fileID,
		SharedWithUserID: recipient.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create share record: %w", err)
	}

	log.Printf("User %d successfully shared file %d with user %d (%s)", ownerID, fileID, recipient.ID, recipientEmail)

	s.auditService.LogActivity(ctx, ownerID, "share:create_user", map[string]interface{}{
		"file_id": fileID,
		"shared_with_user_id": recipient.ID,
		"shared_with_email": recipientEmail,
	})

	return nil
}

func (s *Service) RevokePublicLinks(ctx context.Context, fileID, ownerID int64) error {
	// First, verify ownership to ensure the user can manage this file's shares.
	_, err := s.queries.GetUserFileForDownload(ctx, db.GetUserFileForDownloadParams{ID: fileID, OwnerID: ownerID})
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("file not found or access denied")
		}
		return fmt.Errorf("failed to verify file ownership: %w", err)
	}

	s.auditService.LogActivity(ctx, ownerID, "share:revoke_public", map[string]interface{}{
		"file_id": fileID,
	})

	return s.queries.DeletePublicShareLinksByFileID(ctx, fileID)
}

// UnshareFileWithUser removes a specific user's access to a shared file.
func (s *Service) UnshareFileWithUser(ctx context.Context, fileID, ownerID int64, recipientID int64) error {
	// First, verify ownership of the file.
	_, err := s.queries.GetUserFileForDownload(ctx, db.GetUserFileForDownloadParams{ID: fileID, OwnerID: ownerID})
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("file not found or access denied")
		}
		return fmt.Errorf("failed to verify file ownership: %w", err)
	}

	// --- THIS IS THE FIX ---
	// 1. Perform the database action first.
	err = s.queries.UnshareFileWithUser(ctx, db.UnshareFileWithUserParams{
		UserFileID:       fileID,
		SharedWithUserID: recipientID,
	})
	if err != nil {
		return err // Return the error if the action fails.
	}

	// 2. Only if the action is successful, create the audit log.
	s.auditService.LogActivity(ctx, ownerID, "share:revoke_user", map[string]interface{}{
		"file_id":                 fileID,
		"unshared_from_user_id": recipientID,
	})
	// --- END OF FIX ---

	return nil
}

func (s *Service) GetSharesForFile(ctx context.Context, fileID, ownerID int64) ([]db.GetSharesForFileRow, error) {
	// First, verify ownership.
	_, err := s.queries.GetUserFileForDownload(ctx, db.GetUserFileForDownloadParams{ID: fileID, OwnerID: ownerID})
	if err != nil {
		return nil, fmt.Errorf("file not found or access denied")
	}
	return s.queries.GetSharesForFile(ctx, fileID)
}