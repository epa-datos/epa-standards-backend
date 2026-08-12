// Package api wires the HTTP server: middleware, routes, and dependency
// injection for every resource. main.go only calls api.RunServer().
package api

import (
	"strings"

	"github.com/epa-datos/epa-standards-backend/internal/pkg/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// RunServer builds the gin engine, mounts every route group, and blocks
// serving HTTP traffic on config.Cfg.ServerPort.
func RunServer() {
	cfg := config.Cfg

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(corsMiddleware(cfg.AllowedOrigins))

	registerRoutes(r)

	logrus.Infof("Server running on :%s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		logrus.Fatal(err)
	}
}

func corsMiddleware(allowedOrigins string) gin.HandlerFunc {
	config := cors.DefaultConfig()
	if allowedOrigins != "" {
		config.AllowOrigins = strings.Split(allowedOrigins, ",")
	} else {
		config.AllowAllOrigins = true
	}
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	return cors.New(config)
}
