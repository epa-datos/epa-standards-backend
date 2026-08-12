# Pruebas unitarias

## Cómo correrlas

```bash
go test ./...                 # todos los tests
go test ./... -v              # verboso
go test ./... -cover          # con cobertura por paquete
go test -run TestService_Create ./internal/pkg/service/example   # un test puntual

# o con Makefile
make test
make test-cover
make test-cover-html          # genera coverage.html navegable
```

## Dos niveles de test en este repo

### 1. Servicio (`internal/pkg/service/<recurso>/service_test.go`)

Prueba la **lógica de negocio** en aislamiento, sin HTTP y sin base de
datos: se inyecta un mock de `ports.ExampleRepository` generado por mockery
(ver [MOCKS.md](./MOCKS.md)).

```go
func TestService_Create(t *testing.T) {
    tests := []struct {
        name      string
        example   *entity.Example
        mockSetup func(repo *mocks.ExampleRepository)
        wantErr   error
    }{
        {
            name:    "valid example is created",
            example: &entity.Example{Name: "Sample"},
            mockSetup: func(repo *mocks.ExampleRepository) {
                repo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Example")).Return(nil)
            },
        },
        {
            name:      "empty name is rejected before touching the repository",
            example:   &entity.Example{Name: "   "},
            mockSetup: func(repo *mocks.ExampleRepository) {}, // sin expectativas: no debe llamar al repo
            wantErr:   example.ErrNameRequired,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := mocks.NewExampleRepository(t)
            tt.mockSetup(repo)
            svc := example.NewService(repo)

            err := svc.Create(context.Background(), tt.example)

            if tt.wantErr != nil {
                assert.ErrorIs(t, err, tt.wantErr)
                return
            }
            assert.NoError(t, err)
        })
    }
}
```

Puntos clave:
- **Tabla de casos** (`tests := []struct{...}`) en vez de un test por
  escenario — más fácil de leer y de extender.
- El mock **no** programa expectativas para el caso de validación fallida:
  si el servicio llamara al repositorio por error, el test lo detectaría
  (llamada no esperada → falla).
- `mocks.NewExampleRepository(t)` registra el assert de expectativas al
  final del test automáticamente.

### 2. Handler HTTP (`internal/infrastructure/api/<recurso>/handlers_test.go`)

Prueba la traducción HTTP ↔ servicio usando `net/http/httptest` + un mock de
`ports.ExampleService` (no de repositorio — al handler no le importa cómo
se persisten los datos).

```go
func TestHandler_GetByID(t *testing.T) {
    svc := mocks.NewExampleService(t)
    router := setupRouter(svc) // gin.Engine con solo las rutas de Example

    svc.On("GetByID", mock.Anything, "1").
        Return(&entity.Example{ID: "1", Name: "Sample"}, nil)

    req := httptest.NewRequest(http.MethodGet, "/1", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
}
```

Qué validar en un test de handler:
- Código HTTP correcto para éxito, "no encontrado", error de validación,
  JSON malformado, etc.
- Que los errores de negocio (`errors.Is(err, example.ErrNotFound)`, ...) se
  traduzcan al status HTTP correcto — ver `handleServiceError` en
  `handlers.go`.
- **No** verifiques lógica de negocio aquí (eso va en el test de servicio).

## Cobertura objetivo

| Capa | Objetivo |
|---|---|
| `pkg/service/*` (lógica de negocio) | > 80% |
| `infrastructure/api/*` (handlers) | > 60% |
| `infrastructure/repositories/*` | Pruebas de integración opcionales (requieren DB real); no se exige cobertura de unit test aquí |

## Al agregar un recurso nuevo

1. Copia `internal/pkg/service/example/service_test.go` → tu paquete nuevo,
   renombrando `Example` por tu entidad.
2. Copia `internal/infrastructure/api/example/handlers_test.go` de la misma
   forma.
3. Corre `make mocks` si agregaste interfaces nuevas (ver
   [MOCKS.md](./MOCKS.md)) — los tests no van a compilar hasta que el mock
   exista.
4. `make test` antes de abrir el PR. CI (`.github/workflows/run_tests.yml`)
   corre lo mismo en cada push/PR.
