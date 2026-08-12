package api

import (
	"net/http"

	exampleAPI "github.com/epa-datos/epa-standards-backend/internal/infrastructure/api/example"
	"github.com/epa-datos/epa-standards-backend/internal/infrastructure/api/middlewares"
	"github.com/epa-datos/epa-standards-backend/internal/infrastructure/repositories/postgres"
	exampleSvc "github.com/epa-datos/epa-standards-backend/internal/pkg/service/example"
	"github.com/gin-gonic/gin"
)

// registerRoutes wires every resource's dependencies (repository → service →
// handler) and mounts its routes. This is the single place where you add a
// new resource: create its repo/service/handler packages, then add three
// lines here.
func registerRoutes(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1", middlewares.Auth())

	// --- Example resource ------------------------------------------------
	// Swap postgres.NewExampleRepository(...) for
	// firestore.NewExampleRepository(...) to use Firestore instead — the
	// service and handler below don't change at all.
	exampleRepo := postgres.NewExampleRepository(postgres.NewClient())
	exampleService := exampleSvc.NewService(exampleRepo)
	exampleHandler := exampleAPI.NewHandler(exampleService)
	exampleAPI.RegisterRoutes(v1.Group("/examples"), exampleHandler)
}
