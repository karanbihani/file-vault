package auth

import (
	"context"
	"fmt"
	"time"
	"os"
	"log"
	"strconv"

	"github.com/karanbihani/file-vault/internal/db" // Adjust to your module path
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	queries     *db.Queries
	jwtSecret   []byte
	jwtLifetime time.Duration
}

func NewService(queries *db.Queries) *Service {
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
	// ... (hashing logic is the same)

	// Start a transaction to ensure user creation and role assignment are atomic.
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	qtx := s.queries.WithTx(tx)

	user, err := qtx.CreateUser(ctx, db.CreateUserParams{ /* ... */ })
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Get the 'user' role
	userRole, err := qtx.GetRoleByName(ctx, "user")
	if err != nil {
		return nil, fmt.Errorf("default 'user' role not found: %w", err)
	}

	// Link the new user to the 'user' role.
	if err := qtx.LinkUserToRole(ctx, db.LinkUserToRoleParams{UserID: user.ID, RoleID: userRole.ID}); err != nil {
		return nil, fmt.Errorf("failed to assign role to user: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

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