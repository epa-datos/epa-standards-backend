package example

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/epa-datos/epa-standards-backend/internal/pkg/entity"
	exampleSvc "github.com/epa-datos/epa-standards-backend/internal/pkg/service/example"
	"github.com/epa-datos/epa-standards-backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupRouter wires a bare gin engine with only the Example routes, backed
// by the mocked service — the same pattern to reuse for any other resource.
func setupRouter(svc *mocks.ExampleService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	RegisterRoutes(&r.RouterGroup, NewHandler(svc))
	return r
}

func TestHandler_GetByID(t *testing.T) {
	svc := mocks.NewExampleService(t)
	router := setupRouter(svc)

	t.Run("found", func(t *testing.T) {
		svc.On("GetByID", mock.Anything, "1").Return(&entity.Example{ID: "1", Name: "Sample"}, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var got entity.Example
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
		assert.Equal(t, "Sample", got.Name)
	})

	t.Run("not found", func(t *testing.T) {
		svc.On("GetByID", mock.Anything, "missing").Return(nil, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/missing", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestHandler_Create(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		mockSetup  func(svc *mocks.ExampleService)
		wantStatus int
	}{
		{
			name: "valid payload",
			body: `{"name": "Sample"}`,
			mockSetup: func(svc *mocks.ExampleService) {
				svc.On("Create", mock.Anything, mock.AnythingOfType("*entity.Example")).Return(nil).Once()
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "service rejects empty name",
			body: `{"name": "Sample"}`,
			mockSetup: func(svc *mocks.ExampleService) {
				svc.On("Create", mock.Anything, mock.AnythingOfType("*entity.Example")).Return(exampleSvc.ErrNameRequired).Once()
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "malformed json",
			body:       `{"name":`,
			mockSetup:  func(svc *mocks.ExampleService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := mocks.NewExampleService(t)
			tt.mockSetup(svc)
			router := setupRouter(svc)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
