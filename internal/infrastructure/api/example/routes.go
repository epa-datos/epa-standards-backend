package example

import "github.com/gin-gonic/gin"

// RegisterRoutes mounts the Example resource under the given router group,
// e.g. RegisterRoutes(v1.Group("/examples"), handler).
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.GET("", h.List)
	rg.GET("/:id", h.GetByID)
	rg.POST("", h.Create)
	rg.PUT("/:id", h.Update)
	rg.DELETE("/:id", h.Delete)
}
