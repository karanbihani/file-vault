// ... (The top part of the file is the same) ...
package files

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime"
	"errors"

	"github.com/karanbihani/file-vault/internal/db"      
	"github.com/karanbihani/file-vault/internal/storage" 
	"github.com/gabriel-vasile/mimetype"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	db      *pgxpool.Pool
	queries *db.Queries
	storage *storage.Client
}

func NewService(dbpool *pgxpool.Pool, queries *db.Queries, storageClient *storage.Client) *Service {
	return &Service{
		db:      dbpool,
		queries: queries,
		storage: storageClient,
	}
}

type UploadFileParams struct {
	File        io.Reader
	Filename    string
	ContentType string
	OwnerID     int64
	Description string
	Tags        []string
}

var ErrQuotaExceeded = errors.New("storage quota exceeded")

func (s *Service) UploadFile(ctx context.Context, params UploadFileParams) (*db.UserFile, error) {
	
	user, err := s.queries.GetUserByID(ctx, params.OwnerID)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve user for quota check: %w", err)
	}
	
	var buf bytes.Buffer
	hasher := sha256.New()
	teeReader := io.TeeReader(params.File, hasher)

	size, err := io.Copy(&buf, teeReader)
	if err != nil {
		return nil, fmt.Errorf("could not copy file content to buffer: %w", err)
	}

	if user.StorageUsedBytes+size > user.StorageQuotaBytes {
		return nil, ErrQuotaExceeded
	}

	mtype := mimetype.Detect(buf.Bytes())
	finalMimeType := mtype.String()

	clientBaseMime, _, _ := mime.ParseMediaType(params.ContentType)
	detectedBaseMime, _, _ := mime.ParseMediaType(finalMimeType)

	if clientBaseMime != detectedBaseMime {
		return nil, fmt.Errorf("mime type mismatch: client declared '%s', but content is detected as '%s'", clientBaseMime, detectedBaseMime)
	}

	hash := hex.EncodeToString(hasher.Sum(nil))

	createUserFileParams := db.CreateUserFileParams{
		OwnerID:     params.OwnerID,
		Filename:    params.Filename,
		MimeType:    finalMimeType,
		Description: pgtype.Text{String: params.Description, Valid: params.Description != ""},
		Tags:        params.Tags,
	}

	existingPhysicalFile, err := s.queries.GetPhysicalFileByHash(ctx, hash)
	if err == nil {
		log.Printf("Duplicate file detected. Hash: %s. Incrementing ref count.", hash)
		if _, err := s.queries.IncrementPhysicalFileRefCount(ctx, existingPhysicalFile.ID); err != nil {
			return nil, fmt.Errorf("failed to increment ref count: %w", err)
		}
		
		createUserFileParams.PhysicalFileID = existingPhysicalFile.ID
		newUserFile, err := s.queries.CreateUserFile(ctx, createUserFileParams)
		if err != nil {
			return nil, fmt.Errorf("failed to create user_file for duplicate: %w", err)
		}
		return &newUserFile, nil
	}
	if err != pgx.ErrNoRows {
		return nil, fmt.Errorf("failed to check for existing file: %w", err)
	}

	fileReader := bytes.NewReader(buf.Bytes())
	err = s.storage.Save(ctx, hash, fileReader, size, finalMimeType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to object storage: %w", err)
	}
	log.Printf("Successfully uploaded new file to MinIO. Object name: %s", hash)

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	newPhysicalFile, err := qtx.CreatePhysicalFile(ctx, db.CreatePhysicalFileParams{
		Sha256Hash:  hash,
		SizeBytes:   size,
		StoragePath: hash,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create physical_file: %w", err)
	}

	// Atomically update the user's storage usage since this is a new physical file.
	// CORRECTED: Using the new named parameters from our updated SQL.
	if err := qtx.UpdateUserStorageUsage(ctx, db.UpdateUserStorageUsageParams{
		Amount: size,
		ID:     params.OwnerID,
	}); err != nil {
		return nil, fmt.Errorf("failed to update user storage on upload: %w", err)
	}

	createUserFileParams.PhysicalFileID = newPhysicalFile.ID
	newUserFile, err := qtx.CreateUserFile(ctx, createUserFileParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create user_file: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &newUserFile, nil
}

func (s *Service) ListFiles(ctx context.Context, ownerID int64) ([]db.UserFile, error) {
	return s.queries.ListUserFiles(ctx, ownerID)
}

type DownloadFileResponse struct {
	Data     io.ReadCloser
	Filename string
	Size     int64
}

func (s *Service) DownloadFile(ctx context.Context, fileID, ownerID int64) (*DownloadFileResponse, error) {
	fileMeta, err := s.queries.GetUserFileForDownload(ctx, db.GetUserFileForDownloadParams{ID: fileID, OwnerID: ownerID})
	if err != nil {
		return nil, fmt.Errorf("file not found or access denied: %w", err)
	}

	// CORRECTED: 'object' is now of type *minio.Object, so we can call .Stat()
	object, err := s.storage.Get(ctx, fileMeta.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve file from storage: %w", err)
	}

	stat, err := object.Stat()
	if err != nil {
		return nil, fmt.Errorf("could not get file stats from storage: %w", err)
	}

	return &DownloadFileResponse{
		Data:     object,
		Filename: fileMeta.Filename,
		Size:     stat.Size,
	}, nil
}

func (s *Service) DeleteFile(ctx context.Context, fileID, ownerID int64) error {
	// CORRECTED: Using the new, more secure query that requires both fileID and ownerID.
	fileInfo, err := s.queries.GetFileOwnerAndPhysicalFile(ctx, db.GetFileOwnerAndPhysicalFileParams{ID: fileID, OwnerID: ownerID})
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("file not found or access denied")
		}
		return fmt.Errorf("failed to get file info: %w", err)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	if err := qtx.DeleteUserFile(ctx, fileID); err != nil {
		return fmt.Errorf("failed to delete user file record: %w", err)
	}

	newRefCount, err := qtx.DecrementPhysicalFileRefCount(ctx, fileInfo.PhysicalFileID)
	if err != nil {
		return fmt.Errorf("failed to decrement ref count: %w", err)
	}

	log.Printf("Decremented ref count for physical file ID %d to %d", fileInfo.PhysicalFileID, newRefCount)

	if newRefCount <= 0 {
		log.Printf("Ref count is zero. Deleting physical file and object from storage.")

		// CORRECTED: fileInfo now correctly contains the StoragePath.
		if err := s.storage.Delete(ctx, fileInfo.StoragePath); err != nil {
			return fmt.Errorf("failed to delete object from storage: %w", err)
		}

		if err := qtx.DeletePhysicalFile(ctx, fileInfo.PhysicalFileID); err != nil {
			return fmt.Errorf("failed to delete physical file record: %w", err)
		}

		// CORRECTED: Using the new named parameters from our updated SQL.
		if err := qtx.UpdateUserStorageUsage(ctx, db.UpdateUserStorageUsageParams{
			Amount: -fileInfo.SizeBytes,
			ID:     ownerID,
		}); err != nil {
			return fmt.Errorf("failed to update user storage usage: %w", err)
		}
	}

	return tx.Commit(ctx)
}
