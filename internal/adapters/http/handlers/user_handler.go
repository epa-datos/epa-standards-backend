package handlers

import (
	"encoding/json"
	"net/http"

	"epa-api/internal/domain"
	"epa-api/internal/usecases"
	"epa-api/pkg/logger"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	createUserUsecase *usecases.CreateUserUsecase
	getUserUsecase    *usecases.GetUserUsecase
	logger            *logger.Logger
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(
	createUserUC *usecases.CreateUserUsecase,
	getUserUC *usecases.GetUserUsecase,
	log *logger.Logger,
) *UserHandler {
	return &UserHandler{
		createUserUsecase: createUserUC,
		getUserUsecase:    getUserUC,
		logger:            log,
	}
}

// CreateUser handles POST /api/v1/users
// @Summary Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param request body domain.CreateUserInput true "User data"
// @Success 201 {object} domain.User
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// 1. Parse request
	var req domain.CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		h.logger.Printf("Failed to decode request: %v", err)
		return
	}

	// 2. Call usecase
	user, err := h.createUserUsecase.Execute(r.Context(), req)
	if err != nil {
		// Handle domain errors
		h.handleDomainError(w, err)
		h.logger.Printf("Error creating user: %v", err)
		return
	}

	// 3. Return response
	h.respondWithJSON(w, http.StatusCreated, user)
	h.logger.Printf("User created: %s (%s)", user.ID, user.Email)
}

// GetUser handles GET /api/v1/users/{id}
// @Summary Get user by ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} domain.User
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// 1. Extract ID from path
	id := r.PathValue("id")
	if id == "" {
		h.respondWithError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	// 2. Call usecase
	user, err := h.getUserUsecase.Execute(r.Context(), id)
	if err != nil {
		h.handleDomainError(w, err)
		h.logger.Printf("Error getting user: %v", err)
		return
	}

	// 3. Return response
	h.respondWithJSON(w, http.StatusOK, user)
}

// ErrorResponse is the standard error response structure
type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// respondWithJSON sends a JSON response
func (h *UserHandler) respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// respondWithError sends an error response
func (h *UserHandler) respondWithError(w http.ResponseWriter, statusCode int, message string) {
	h.respondWithJSON(w, statusCode, ErrorResponse{
		Error: message,
		Code:  statusCode,
	})
}

// handleDomainError converts domain errors to HTTP responses
func (h *UserHandler) handleDomainError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrUserNotFound:
		h.respondWithError(w, http.StatusNotFound, "User not found")
	case domain.ErrUserAlreadyExists:
		h.respondWithError(w, http.StatusConflict, "User with this email already exists")
	case domain.ErrInvalidEmail:
		h.respondWithError(w, http.StatusBadRequest, "Invalid email format")
	case domain.ErrInvalidName:
		h.respondWithError(w, http.StatusBadRequest, "Name is required and must be at least 2 characters")
	default:
		h.respondWithError(w, http.StatusInternalServerError, "Internal server error")
	}
}
