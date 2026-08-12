// Package middlewares holds gin middleware shared across route groups.
package middlewares

import (
	"net/http"
	"strings"

	"github.com/epa-datos/epa-standards-backend/internal/pkg/config"
	"github.com/gin-gonic/gin"
)

// Auth is a minimal placeholder middleware: it checks for a static bearer
// token from config.Cfg.APIAuthToken. It exists so the folder/pattern is in
// place — replace its body with real auth (Firebase ID tokens, JWT, OAuth2,
// an API gateway check, ...) before using this in production.
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.Cfg.APIAuthToken == "" {
			// No token configured (e.g. local dev): skip auth entirely.
			c.Next()
			return
		}

		header := c.GetHeader("Authorization")
		token := strings.TrimPrefix(header, "Bearer ")
		if token == "" || token != config.Cfg.APIAuthToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or missing token"})
			return
		}

		c.Next()
	}
}
