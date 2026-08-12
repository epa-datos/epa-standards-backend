// Package example is the reference vertical slice for the HTTP layer: one
// handlers.go + routes.go pair per resource. Copy this folder (and its
// _test.go) whenever you add a new resource to the API.
package example

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/epa-datos/epa-standards-backend/internal/pkg/entity"
	"github.com/epa-datos/epa-standards-backend/internal/pkg/ports"
	exampleSvc "github.com/epa-datos/epa-standards-backend/internal/pkg/service/example"
	"github.com/epa-datos/epa-standards-backend/internal/pkg/utils"
	"github.com/gin-gonic/gin"
)

// Handler groups the HTTP handlers for the Example resource. It only
// depends on ports.ExampleService (an interface), so tests can inject
// mocks.ExampleService instead of a real service.
type Handler struct {
	svc ports.ExampleService
}

// NewHandler builds a Handler for the given service implementation.
func NewHandler(svc ports.ExampleService) *Handler {
	return &Handler{svc: svc}
}

// List handles GET /api/v1/examples
func (h *Handler) List(c *gin.Context) {
	rawOffset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)
	rawLimit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)
	offset, limit := utils.ParsePage(rawOffset, rawLimit)

	resp, err := h.svc.List(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetByID handles GET /api/v1/examples/:id
func (h *Handler) GetByID(c *gin.Context) {
	item, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "example not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// Create handles POST /api/v1/examples
func (h *Handler) Create(c *gin.Context) {
	var req entity.Example
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		h.handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, req)
}

// Update handles PUT /api/v1/examples/:id
func (h *Handler) Update(c *gin.Context) {
	var req entity.Example
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.Update(c.Request.Context(), c.Param("id"), &req); err != nil {
		h.handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, req)
}

// Delete handles DELETE /api/v1/examples/:id
func (h *Handler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// handleServiceError maps known sentinel errors from the service layer to
// HTTP status codes. This is the pattern to extend as you add more errors.
func (h *Handler) handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, exampleSvc.ErrNameRequired):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, exampleSvc.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
