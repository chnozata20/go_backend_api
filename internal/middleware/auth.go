package middleware

import (
	"net/http"
	"strings"

	"github.com/cihan-ozata/backend-path/internal/domain"
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

        // TODO: JWT token doğrulama işlemi burada yapılacak
        // Şimdilik basit bir kontrol yapıyoruz
        token := strings.TrimPrefix(authHeader, "Bearer ")
        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // TODO: Token'dan kullanıcı bilgilerini çek
        // Şimdilik örnek bir kullanıcı oluşturuyoruz
        user := &domain.User{
            ID:    1,
            Role: domain.UserRoleUser,
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

        currentUser := user.(*domain.User)
        if currentUser.Role != requiredRole {
            c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
            c.Abort()
            return
        }

        c.Next()
    }
} 