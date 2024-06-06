package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Kachyr/findyourpet/findyourpet-backend/internal/store/users"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// RequireAuth is a middleware that checks for a valid JWT token in the Authorization cookie.
func RequireAuth(userStore users.UserStoreI) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the cookie
		tokenString, err := c.Cookie("Authorization")
		if err != nil {
			log.Error().Err(err).Msg("Failed to get Authorization cookie")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Validate the token
		token, err := validateToken(tokenString)
		if err != nil {
			log.Error().Err(err).Msg("Invalid token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Check if the token is expired
		if isExpired, err := isTokenExpired(token); err != nil || isExpired {
			log.Error().Err(err).Bool("expired", isExpired).Msg("Token is expired or could not be checked")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Extract the user ID from the token claims
		userID, err := extractUserID(token)
		if err != nil {
			log.Error().Err(err).Msg("Failed to extract user ID from token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Retrieve the user from the store
		user, err := userStore.GetByID(userID)
		if err != nil {
			log.Error().Err(err).Str("userID", userID.String()).Msg("User not found")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Set the user in the context for later use
		c.Set("user", user)

		// Proceed to the next handler
		c.Next()
	}
}

// validateToken validates the JWT token and returns the parsed token.
func validateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	return token, err
}

// isTokenExpired checks if the token is expired.
func isTokenExpired(token *jwt.Token) (bool, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return true, fmt.Errorf("invalid token claims")
	}

	expTime, ok := claims["exp"].(float64)
	if !ok {
		return true, fmt.Errorf("invalid exp claim in token")
	}

	return float64(time.Now().Unix()) >= expTime, nil
}

// extractUserID extracts the user ID from the token claims.
func extractUserID(token *jwt.Token) (uuid.UUID, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid token claims")
	}

	userIDstr, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid sub claim in token")
	}

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}
