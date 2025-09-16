package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// --- Database Connection ---
	// We read the database connection string from an environment variable.
	// This is a best practice for security and flexibility.
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// pgxpool is a high-performance connection pool for PostgreSQL.
	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		// Use log.Fatalf to exit the application if the DB connection fails.
		// The application cannot run without the database.
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	// defer dbpool.Close() ensures that the connection pool is closed when main() exits.
	defer dbpool.Close()

	fmt.Println("Successfully connected to the database!")

	// --- Gin Web Server Setup ---
	router := gin.Default()

	// A simple health check endpoint to verify that the server is running.
	router.GET("/health", func(c *gin.Context) {
		// We can also add a database ping here to check DB health.
		err := dbpool.Ping(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "DOWN", "error": "database connection error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	fmt.Println("Starting server on port 8080...")
	// router.Run() starts the HTTP server and listens for incoming requests.
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}