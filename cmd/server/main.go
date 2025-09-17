package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/karanbihani/file-vault/internal/api" // IMPORTANT: Adjust this import path
	"github.com/karanbihani/file-vault/internal/core/files"
)

func main() {
	// ... (database connection and service initialization are the same)
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

	fileService := files.NewService(dbpool, "./uploads")
	log.Println("File service initialized.")

	// --- Gin Web Server Setup ---
	// UPDATED: We now pass the fileService into our router setup function.
	router := api.SetupRouter(dbpool, fileService)

	log.Println("Starting server on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}