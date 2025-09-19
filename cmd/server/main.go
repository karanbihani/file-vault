package main

import (
	"context"
	"log"
	"os"

	"github.com/karanbihani/file-vault/internal/api"      // Adjust path
	"github.com/karanbihani/file-vault/internal/auth"     // Adjust path
	"github.com/karanbihani/file-vault/internal/core/files" // Adjust path
	"github.com/karanbihani/file-vault/internal/db"       // Add this import
	"github.com/karanbihani/file-vault/internal/storage"  // Adjust path
	"github.com/karanbihani/file-vault/internal/core/stats"
	"github.com/karanbihani/file-vault/internal/core/shares"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// --- Database Connection ---
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}
	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()
	log.Println("Successfully connected to the database!")

	// --- SQLC Querier Initialization ---
	// We create the querier object ONCE here.
	queries := db.New(dbpool)
	

	// --- MinIO Client Initialization ---
	minioConfig := storage.Config{
		Endpoint:        os.Getenv("MINIO_ENDPOINT"),
		AccessKeyID:     os.Getenv("MINIO_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("MINIO_SECRET_ACCESS_KEY"),
		BucketName:      os.Getenv("MINIO_BUCKET_NAME"),
		UseSSL:          false, // For local development
	}
	storageClient := storage.NewClient(context.Background(), minioConfig)
	log.Println("MinIO client initialized and bucket is ready.")

	// --- Initialize Services ---
	// We inject the shared 'queries' object into both services.
	authService := auth.NewService(queries)
	fileService := files.NewService(dbpool, queries, storageClient)
	sharesService := shares.NewService(queries, storageClient) // Create the shares service
	statsService := stats.NewService(queries)
	log.Println("Services initialized.")

	// --- Gin Web Server Setup ---
	router := api.SetupRouter(dbpool, fileService, authService, sharesService, statsService)

	log.Println("Starting server on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}