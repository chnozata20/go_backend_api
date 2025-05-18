package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
    validate = validator.New()
    // Custom validations
    validate.RegisterValidation("password", validatePassword)
    validate.RegisterValidation("currency", validateCurrency)
}

// validatePassword checks if password meets requirements
func validatePassword(fl validator.FieldLevel) bool {
    password := fl.Field().String()
    if len(password) < 8 {
        return false
    }
    hasUpper := strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
    hasLower := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
    hasNumber := strings.ContainsAny(password, "0123456789")
    hasSpecial := strings.ContainsAny(password, "!@#$%^&*()_+-=[]{}|;:,.<>?")
    return hasUpper && hasLower && hasNumber && hasSpecial
}

// validateCurrency checks if currency code is valid
func validateCurrency(fl validator.FieldLevel) bool {
    currency := fl.Field().String()
    validCurrencies := map[string]bool{
        "USD": true,
        "EUR": true,
        "TRY": true,
        "GBP": true,
    }
    return validCurrencies[currency]
}

// ValidationMiddleware handles request validation
func ValidationMiddleware(schema interface{}) gin.HandlerFunc {
    return func(c *gin.Context) {
        if err := c.ShouldBindJSON(schema); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "Invalid request format",
                "details": err.Error(),
            })
            c.Abort()
            return
        }

        if err := validate.Struct(schema); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "Validation failed",
                "details": err.Error(),
            })
            c.Abort()
            return
        }

        c.Set("validated_data", schema)
        c.Next()
    }
} 