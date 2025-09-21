package api

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// getJWTSecret retrieves the JWT secret from an environment variable.
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		panic("JWT_SECRET environment variable is not set")
	}
	return []byte(secret)
}

var jwtSecret = getJWTSecret()

// AuthMiddleware creates a Gin middleware for JWT authentication.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header from the request.
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		// The header should be in the format "Bearer <token>". We split it to get the token part.
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}
		tokenString := parts[1]

		// Now, we parse and validate the token.
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is what we expect (HS256).
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token: " + err.Error()})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// The token is valid. We extract the user ID ("sub" claim) from the claims.
			userIDFloat, ok := claims["sub"].(float64)
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid subject claim in token"})
				return
			}
			userID := int64(userIDFloat)

			// CRITICAL: We add the authenticated user's ID to the Gin context.
			// This is how our downstream handlers will know who the user is.
			c.Set("userID", userID)

			// c.Next() passes control to the next handler in the chain.
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		}
	}
}

type client struct {
	lastSeen time.Time
	requests int
}

var (
	clients = make(map[string]*client)
	mu      sync.Mutex
)

// RateLimiter creates a Gin middleware for simple IP-based rate limiting.
func RateLimiter(limit int, window time.Duration) gin.HandlerFunc {
	// Start a background goroutine to clean up old clients periodically.
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, c := range clients {
				if time.Since(c.lastSeen) > window {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		mu.Lock()
		defer mu.Unlock()

		ip := c.ClientIP()

		// If client is not in the map, add them.
		if _, found := clients[ip]; !found {
			clients[ip] = &client{lastSeen: time.Now(), requests: 1}
			c.Next()
			return
		}

		// If client is in the map, check their request time and count.
		c_ := clients[ip]
		if time.Since(c_.lastSeen) > window {
			c_.lastSeen = time.Now()
			c_.requests = 1
			c.Next()
			return
		}

		c_.requests++
		if c_.requests > limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}

		c.Next()
	}
}