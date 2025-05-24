package middleware

import (
	"net/http"
	"time"

	"go-crime_map_backend/internal/domain/repositories"

	"github.com/gin-gonic/gin"
)

// APIKeyHeader es el nombre del header que contiene la API Key
const APIKeyHeader = "X-API-Key"

// AuthMiddleware crea un middleware de autenticación
func AuthMiddleware(securityRepo repositories.SecurityRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader(APIKeyHeader)
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API Key requerida"})
			c.Abort()
			return
		}

		key, err := securityRepo.ValidateAPIKey(c.Request.Context(), apiKey)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "API Key inválida"})
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if key == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API Key inválida"})
			c.Abort()
			return
		}

		if key != nil {
			if !key.IsActiveBool {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "API Key inactiva"})
				c.Abort()
				return
			}
		}

		if time.Now().After(key.ExpiresAt) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API Key expirada"})
			c.Abort()
			return
		}

		c.Set("user_id", key.ID)
		c.Next()
	}
}
