package auth

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/karanbihani/file-vault/internal/db" // Adjust to your module path
)

type Service struct {
	db          *pgxpool.Pool 
	queries     *db.Queries
	jwtSecret   []byte
	jwtLifetime time.Duration
}

func NewService(dbpool *pgxpool.Pool, queries *db.Queries) *Service {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		log.Fatal("JWT_SECRET_KEY environment variable is not set")
	}

	lifetimeHoursStr := os.Getenv("JWT_LIFETIME_HOURS")
	if lifetimeHoursStr == "" {
		lifetimeHoursStr = "24" // Default to 24 hours
	}
	lifetimeHours, err := strconv.Atoi(lifetimeHoursStr)
	if err != nil {
		log.Fatalf("Invalid JWT_LIFETIME_HOURS: %v", err)
	}

	return &Service{
		db:          dbpool, 
		queries:     queries,
		jwtSecret:   []byte(secret),
		jwtLifetime: time.Hour * time.Duration(lifetimeHours),
	}
}

type RegisterUserParams struct {
	Email    string
	Password string
}

// RegisterUser creates a new user, hashes their password, and saves it to the database.
func (s *Service) RegisterUser(ctx context.Context, params RegisterUserParams) (*db.User, error) {
	log.Println("Starting user registration process")
	log.Printf("Input params: %+v", params)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Start a transaction to ensure user creation and role assignment are atomic.
	tx, err := s.db.Begin(ctx)
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return nil, fmt.Errorf("could not begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			log.Println("Rolling back transaction due to error")
			tx.Rollback(ctx)
		}
	}()

	qtx := s.queries.WithTx(tx)

	log.Println("Creating user in the database")
	user, err := qtx.CreateUser(ctx, db.CreateUserParams{
		Email:        params.Email,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	log.Println("Fetching default 'user' role")
	userRole, err := qtx.GetRoleByName(ctx, "user")
	if err != nil {
		log.Printf("Error fetching 'user' role: %v", err)
		return nil, fmt.Errorf("default 'user' role not found: %w", err)
	}

	log.Println("Linking user to 'user' role")
	if err := qtx.LinkUserToRole(ctx, db.LinkUserToRoleParams{UserID: user.ID, RoleID: userRole.ID}); err != nil {
		log.Printf("Error linking user to role: %v", err)
		return nil, fmt.Errorf("failed to assign role to user: %w", err)
	}

	log.Println("Committing transaction")
	if err := tx.Commit(ctx); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Println("User registration successful")
	return &user, nil
}

type LoginUserParams struct {
	Email    string
	Password string
}

func (s *Service) LoginUser(ctx context.Context, params LoginUserParams) (string, error) {
	user, err := s.queries.GetUserByEmail(ctx, params.Email)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(params.Password))
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"iat": time.Now().Unix(),
		// Use the configured lifetime from the service.
		"exp": time.Now().Add(s.jwtLifetime).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Use the configured secret from the service.
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}

	return tokenString, nil
}