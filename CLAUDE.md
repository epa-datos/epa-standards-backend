# CLAUDE.md - EPA Digital Standard Backend API

Este es un template estándar de API Go con arquitectura hexagonal. Léelo completamente antes de empezar.

## 🚀 Comienza Aquí (5 min)

```bash
# 1. Descarga dependencias
go mod download

# 2. Copia variables de entorno
cp .env.example .env

# 3. Inicia el servidor
go run ./cmd/api

# Deberías ver: "Server running on :8080"
# API disponible en: http://localhost:8080
```

---

## 📋 Comandos Comunes

```bash
# Desarrollo
go run ./cmd/api              # Inicia servidor localmente
go run ./cmd/api/main.go      # Alternativa

# Testing
go test ./...                 # Todos los tests
go test ./... -v              # Tests verbosos
go test ./... -cover          # Con cobertura
go test -run TestNombreTest   # Test específico

# Linting
go vet ./...                  # Go vet check
golangci-lint run             # Linting completo (si tienes instalado)

# Building
go build -o api ./cmd/api    # Compilar binario
make build                    # Si tienes Makefile

# Docker
docker build -t api:latest .
docker run -p 8080:8080 api:latest

# OpenAPI/Swagger
swag init -g cmd/api/main.go # Generar docs/swagger.json
# Luego: http://localhost:8080/swagger/index.html
```

---

## 🏗️ Arquitectura: Hexagonal (Ports & Adapters)

**Principio:** La lógica de negocio es **independiente** de HTTP, DB, o frameworks.

```
┌─────────────────────────────────────────┐
│           HTTP Handlers (Adapters)       │  ← Reciben requests
├─────────────────────────────────────────┤
│           Use Cases (Orchestration)      │  ← Orquestan domain + puertos
├─────────────────────────────────────────┤
│        Domain (Business Logic)           │  ← Pure Go, sin deps externas
├─────────────────────────────────────────┤
│     Repositories & Services (Adapters)   │  ← Implementan interfaces
├─────────────────────────────────────────┤
│     External APIs, Database, Caches      │  ← Infraestructura
└─────────────────────────────────────────┘
```

**Flujo real: Crear Usuario**

```
1. HTTP POST /api/v1/users {"email": "..."}
2. Handler → handler.CreateUser() ← HTTP adapter
3. Handler → usecase.CreateUser(ctx, req) ← Orquestación
4. Usecase → userService.Validate() ← Domain interface
5. Usecase → userRepo.Save() ← Adapter interface (se inyecta)
6. Repository → database.Insert() ← Infraestructura
7. Response fluye de vuelta con 201 Created
```

### Estructura de Directorios

```
cmd/api/
├── main.go              # Entry point, dependency injection
└── config.go            # Configuration loading (optional)

internal/
├── domain/              # 🔴 SIN dependencias externas
│   ├── entities.go      # User, Product, etc. (tipos puros)
│   ├── errors.go        # DomainError, código de error
│   └── interfaces.go    # Interfaces que adapters implementan
│                         # ej: UserRepository, EmailService
│
├── usecases/            # 🟡 Orquestación de lógica
│   ├── create_user.go   # CreateUserUsecase
│   ├── get_user.go
│   └── delete_user.go
│
└── adapters/            # 🔵 Implementaciones de infraestructura
    ├── http/
    │   ├── handlers/
    │   │   ├── user_handler.go    # HTTP handlers
    │   │   └── product_handler.go
    │   ├── middleware/
    │   │   ├── auth.go
    │   │   ├── cors.go
    │   │   └── error_handler.go
    │   ├── router.go              # Routes registration
    │   └── response.go            # Response helpers
    │
    ├── persistence/
    │   ├── user_repository.go     # Implementa domain.UserRepository
    │   └── product_repository.go
    │
    └── external/
        ├── email_service.go       # Implementa domain.EmailService
        └── payment_service.go

pkg/                    # Utilidades compartidas
├── logger/             # Logger utilities
├── errors/             # Error types
├── validator/          # Validation helpers
└── middleware/         # Shared middleware

docs/
├── swagger.json        # ⚠️ AUTO-GENERADO (no editar)
└── api_design.md       # Notas de diseño

.github/workflows/
├── test.yml            # Run tests on PR
├── openapi-sync.yml    # Sync to Postman
└── deploy.yml          # Deploy to Cloud Run
```

### Reglas de Oro

**✅ DO:**
- Domain depende solo de Go estándar
- Handlers inyectan dependencias (usecase)
- Usecases inyectan dependencias (repositories, servicios)
- Repositories implementan interfaces del domain

**❌ DON'T:**
- Handler llama directo a DB (incorrecto: Handler → Database)
- Domain importa adapters (imports circulares)
- Usecase conoce detalles HTTP (status codes, headers)

---

## 🔌 Handlers & HTTP

### Handler Pattern

```go
// internal/adapters/http/handlers/user_handler.go

// UserHandler contiene dependencias inyectadas
type UserHandler struct {
    createUserUsecase *usecases.CreateUserUsecase
    getUserUsecase    *usecases.GetUserUsecase
    logger            *log.Logger
}

// Constructor
func NewUserHandler(
    createUserUC *usecases.CreateUserUsecase,
    getUserUC *usecases.GetUserUsecase,
    logger *log.Logger,
) *UserHandler {
    return &UserHandler{
        createUserUsecase: createUserUC,
        getUserUsecase:    getUserUC,
        logger:            logger,
    }
}

// POST /api/v1/users
// @Summary Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "User data"
// @Success 201 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // 1. Parse request
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
        return
    }
    
    // 2. Call usecase (no lógica aquí, solo delegation)
    user, err := h.createUserUsecase.Execute(r.Context(), req)
    if err != nil {
        // 3. Handle domain errors
        h.handleError(w, err)
        return
    }
    
    // 4. Return response
    h.respondWithJSON(w, http.StatusCreated, user)
}

// Helper para respuestas
func (h *UserHandler) respondWithJSON(w http.ResponseWriter, code int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(data)
}

func (h *UserHandler) respondWithError(w http.ResponseWriter, code int, msg string) {
    h.respondWithJSON(w, code, map[string]string{"error": msg})
}
```

### OpenAPI Comments (IMPORTANTE)

Cada handler debe tener comentarios OpenAPI. Ejemplo:

```go
// GET /api/v1/users/{id}
// @Summary Get user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} UserResponse
// @Failure 404 {object} ErrorResponse "User not found"
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    // ...
}
```

**Después de agregar/cambiar handlers:**
```bash
swag init -g cmd/api/main.go
git add docs/swagger.json
```

---

## 🧪 Testing

### Estructura

```
internal/
├── domain/
│   └── entities_test.go       # Tests de entities
├── usecases/
│   └── create_user_test.go    # Tests de usecase (mocks)
└── adapters/
    └── http/
        └── handlers/
            └── user_handler_test.go  # Tests de handler
```

### Ejemplo: Test de Usecase

```go
// internal/usecases/create_user_test.go

func TestCreateUser(t *testing.T) {
    tests := []struct {
        name      string
        req       CreateUserRequest
        mockRepo  *mockUserRepository
        wantError bool
        wantUser  *domain.User
    }{
        {
            name: "valid user creation",
            req: CreateUserRequest{
                Email: "user@example.com",
                Name:  "John Doe",
            },
            mockRepo:  &mockUserRepository{},
            wantError: false,
            wantUser: &domain.User{
                Email: "user@example.com",
                Name:  "John Doe",
            },
        },
        {
            name: "invalid email",
            req: CreateUserRequest{
                Email: "invalid",
                Name:  "John",
            },
            wantError: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            usecase := usecases.NewCreateUserUsecase(tt.mockRepo)
            
            user, err := usecase.Execute(context.Background(), tt.req)
            
            if (err != nil) != tt.wantError {
                t.Errorf("got error %v, want %v", err != nil, tt.wantError)
            }
            
            if !tt.wantError && user.Email != tt.wantUser.Email {
                t.Errorf("got email %s, want %s", user.Email, tt.wantUser.Email)
            }
        })
    }
}

// Mock para testing
type mockUserRepository struct {
    SaveFunc func(ctx context.Context, user *domain.User) error
}

func (m *mockUserRepository) Save(ctx context.Context, user *domain.User) error {
    return m.SaveFunc(ctx, user)
}
```

### Coverage Target
- Domain & Usecases: >85%
- Handlers: >70%
- (Helpers y test utilities: pueden ser menores)

---

## 🐳 Docker & Cloud Run

### Dockerfile (Multi-stage)

```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o api ./cmd/api

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /build/api .
EXPOSE 8080
HEALTHCHECK --interval=30s CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1
CMD ["./api"]
```

### Deploy a Cloud Run

```bash
# Build locally
docker build -t gcr.io/PROJECT/api:latest .

# Push
docker push gcr.io/PROJECT/api:latest

# Deploy
gcloud run deploy api \
  --image gcr.io/PROJECT/api:latest \
  --region us-central1 \
  --set-env-vars GCP_PROJECT_ID=PROJECT \
  --allow-unauthenticated
```

GitHub Actions hace esto automáticamente en releases.

---

## 📝 Ejemplos en Este Repo

- `internal/domain/user.go` - Entity de ejemplo
- `internal/usecases/create_user.go` - Usecase de ejemplo
- `internal/adapters/http/handlers/user_handler.go` - Handler de ejemplo
- `internal/adapters/persistence/user_repository.go` - Repository de ejemplo
- `cmd/api/main.go` - Dependency injection setup

**Cómo usar:** Copia estos archivos como plantilla para tus propias features.

---

## 🔐 Variables de Entorno

Copia `.env.example` a `.env`:

```bash
cp .env.example .env
```

**En `.env`:**
```env
PORT=8080
GCP_PROJECT_ID=your-project-id
LOG_LEVEL=debug
# Agrega más según necesites
```

**Nunca commites `.env`** (está en .gitignore)

En `cmd/api/main.go`:
```go
func main() {
    cfg := config.Load() // Lee de .env
    log.Printf("Starting server on port %s", cfg.Port)
    // ...
}
```

---

## 🔄 Git Workflow

### Branch Naming
```
feature/user-authentication    # Nueva feature
fix/email-validation-bug       # Bug fix
refactor/error-handling        # Mejora de código
```

### PR Process
1. Crea rama desde `staging`
2. Haz commit significativos
3. Push a remote
4. Abre PR a `staging`
5. Espera CI (tests, linting)
6. Get 1 approval
7. Merge

### Commits
```bash
git commit -m "feat: add user authentication endpoint"
git commit -m "fix: validate email format"
git commit -m "test: add user creation tests"
git commit -m "docs: update API spec"
```

---

## 💡 Errores Comunes & Soluciones

**Error: "module not found"**
```bash
go mod tidy
go mod download
```

**Error: "connection refused" en tests**
→ Estás intentando usar DB real. Usa mocks en tests.

**Handler devuelve 500**
→ Revisa logs. Usa `h.logger.Printf()` para debug.

**OpenAPI spec desactualizado**
```bash
swag init -g cmd/api/main.go
git add docs/swagger.json
```

---

## 📚 Recursos

- **Go Docs:** https://golang.org/doc/
- **Hexagonal Architecture:** https://alistair.cockburn.us/hexagonal-architecture/
- **OpenAPI 3.0:** https://swagger.io/specification/
- **Cloud Run:** https://cloud.google.com/run/docs

---

## 🎯 Siguientes Pasos

1. ✅ Lee esta guía completa
2. ✅ Ejecuta `go run ./cmd/api` y verifica que corre
3. ✅ Revisa los archivos de ejemplo en `internal/`
4. ✅ Abre `http://localhost:8080/swagger/index.html` (si swag está correctamente configurado)
5. ✅ Crea tu primer handler copiando `user_handler.go`
6. ✅ Escribe tests para tu handler
7. ✅ Haz un PR a `staging`

---

## 🤖 Skills Disponibles

Estos skills de Claude están disponibles para ayudarte con proyectos Go:

### `go-api-scaffold`
Crea un nuevo repositorio Go API basado en esta plantilla estándar.

**Cuándo usar:** Al iniciar un nuevo servicio Go para EPA Digital
```
Claude: "Create a new Go API for the billing service"
```

### `go-client-scaffold`
Crea una librería cliente Go reutilizable para consumir APIs.

**Cuándo usar:** Necesitas un cliente Go para otro servicio
```
Claude: "Create a Go client for the analytics service"
```

### `generate-openapi`
Genera o actualiza documentación OpenAPI/Swagger para tu API.

**Cuándo usar:** Después de agregar nuevos handlers
```
Claude: "Update OpenAPI docs for my API"
```

### `validate-pr-format`
Verifica que tu PR siga los estándares EPA Digital antes de submitear.

**Cuándo usar:** Antes de abrir un PR
```
Claude: "Validate my PR format"
```

### `git-flow-guide`
Guía interactiva para branching strategy y workflow de EPA Digital.

**Cuándo usar:** Duda sobre qué rama crear
```
Claude: "What branch should I create for a feature?"
```

---

## ❓ ¿Preguntas?

- Revisa ejemplos en `internal/`
- Revisa tests en `*_test.go`
- Lee comentarios en el código
- Usa los skills disponibles (ver arriba)
- Pregunta en el equipo

**Buena suerte!** 🚀
