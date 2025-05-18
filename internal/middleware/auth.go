package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/cihan-ozata/backend-path/internal/domain"
	"github.com/cihan-ozata/backend-path/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware handles JWT authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, err := jwt.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		user := &domain.User{
			ID:       claims.UserID,
			Username: claims.Username,
			Email:    claims.Email,
			Role:     domain.UserRole(claims.Role),
		}

		c.Set("user", user)
		c.Next()
	}
}

// RoleMiddleware checks if user has required role
func RoleMiddleware(requiredRole domain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}
		log.Printf("DEBUG: currentUser.Role = %s, requiredRole = %s\n", user.(*domain.User).Role, requiredRole)

		currentUser := user.(*domain.User)
		if currentUser.Role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
} 