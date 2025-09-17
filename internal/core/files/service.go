package files

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime"
	"os"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/karanbihani/file-vault/internal/db"
)

// Service struct holds the dependencies for the file service,
// primarily the database connection pool and the generated queries.
type Service struct {
	db          *pgxpool.Pool
	queries     *db.Queries
	storageRoot string
}

// NewService is a factory function that creates a new File Service.
func NewService(dbpool *pgxpool.Pool, storageRoot string) *Service {
	if err := os.MkdirAll(storageRoot, 0755); err != nil {
		log.Fatalf("Could not create storage directory: %v", err)
	}
	log.Printf("Storage root is set to: %s", storageRoot)
	return &Service{
		db:          dbpool,
		queries:     db.New(dbpool),
		storageRoot: storageRoot,
	}
}

type UploadFileParams struct {
	File        io.Reader
	Filename    string
	ContentType string
	OwnerID     int64
	Description string   // ADDED
	Tags        []string // ADDED
}


func (s *Service) UploadFile(ctx context.Context, params UploadFileParams) (*db.UserFile, error) {
	
	tempFile, err := os.CreateTemp(s.storageRoot, "upload-*.tmp")
	if err != nil {
		return nil, fmt.Errorf("could not create temp file: %w", err)
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	hasher := sha256.New()
	teeReader := io.TeeReader(params.File, hasher)

	size, err := io.Copy(tempFile, teeReader)
	if err != nil {
		return nil, fmt.Errorf("could not copy file content: %w", err)
	}

	// Validate the actual file content against the declared MIME type.
	if _, err := tempFile.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("could not seek temp file: %w", err)
	}

	// 1. Detect the full media type from the file content.
	detectedMtype, err := mimetype.DetectReader(tempFile)
	if err != nil {
		return nil, fmt.Errorf("could not detect mime type: %w", err)
	}

	// 2. Parse BOTH the client-declared and server-detected types to get their base type.
	// This correctly handles cases like 'text/plain' vs 'text/plain; charset=utf-8'.
	clientBaseMime, _, err := mime.ParseMediaType(params.ContentType)
	if err != nil {
		return nil, fmt.Errorf("could not parse client content type: %w", err)
	}

	detectedBaseMime, _, err := mime.ParseMediaType(detectedMtype.String())
	if err != nil {
		return nil, fmt.Errorf("could not parse detected mime type: %w", err)
	}
	
	// 3. Compare the base types.
	if clientBaseMime != detectedBaseMime {
		return nil, fmt.Errorf("mime type mismatch: client declared '%s', but content is detected as '%s'", clientBaseMime, detectedBaseMime)
	}
	// Use the more specific, server-detected MIME type for storage.
	actualMimeType := detectedMtype.String()
	// --- END CORRECTION ---

	// Compute the SHA-256 hash of the file content.
	hash := hex.EncodeToString(hasher.Sum(nil))

	existingPhysicalFile, err := s.queries.GetPhysicalFileByHash(ctx, hash)
	if err == nil {
		log.Printf("Duplicate file detected. Hash: %s. Incrementing ref count for physical file ID: %d", hash, existingPhysicalFile.ID)
		if _, err := s.queries.IncrementPhysicalFileRefCount(ctx, existingPhysicalFile.ID); err != nil {
			return nil, fmt.Errorf("failed to increment ref count: %w", err)
		}
		
		// Create the user_file record pointing to the existing physical file.
		newUserFile, err := s.queries.CreateUserFile(ctx, db.CreateUserFileParams{
			OwnerID:        params.OwnerID,
			PhysicalFileID: existingPhysicalFile.ID,
			Filename:       params.Filename,
			MimeType:       actualMimeType,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create user_file for duplicate: %w", err)
		}
		return &newUserFile, nil
	}

	// Name unique file by its hash to ensure uniqueness.
	finalPath := filepath.Join(s.storageRoot, hash)
	if err := os.Rename(tempFile.Name(), finalPath); err != nil {
		return nil, fmt.Errorf("could not move temp file to storage: %w", err)
	}

	// 7. Create the physical_file and user_file records in a single database transaction.
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // Rollback is a no-op if the transaction is committed.

	qtx := s.queries.WithTx(tx)

	newPhysicalFile, err := qtx.CreatePhysicalFile(ctx, db.CreatePhysicalFileParams{
		Sha256Hash:  hash,
		SizeBytes:   size,
		StoragePath: finalPath,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create physical_file: %w", err)
	}

	newUserFile, err := qtx.CreateUserFile(ctx, db.CreateUserFileParams{
		OwnerID:        params.OwnerID,
		PhysicalFileID: newPhysicalFile.ID,
		Filename:       params.Filename,
		MimeType:       actualMimeType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user_file: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("File is stored at: %s", finalPath)
	
	return &newUserFile, nil
}
