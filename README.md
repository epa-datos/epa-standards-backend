# EPA Digital - Standard Backend API

Template de API Go con arquitectura hexagonal, lista para clonar y empezar a desarrollar.

**Usa este repo como punto de partida para cualquier API Go en EPA Digital.**

## 🚀 Quick Start

```bash
# 1. Clone
git clone https://github.com/epadigital/epa-standards-backend.git my-api
cd my-api

# 2. Setup
cp .env.example .env
go mod download

# 3. Run
go run ./cmd/api

# 4. Test
go test ./...

# 5. Read CLAUDE.md
cat CLAUDE.md  # ← Leelo completamente
```

## 📚 Documentación

**Antes de empezar, lee:**
- **[CLAUDE.md](./CLAUDE.md)** - Guía completa de arquitectura y desarrollo
- **[CONTRIBUTING.md](../epa-standards/CONTRIBUTING.md)** - Cómo contribuir
- **[BRANCHING-STRATEGY.md](../epa-standards/docs/BRANCHING-STRATEGY.md)** - Git workflow

## 📁 Estructura

```
cmd/api/                    # Entry point
internal/
  ├── domain/              # Business logic (entities, interfaces)
  ├── usecases/            # Orchestration (application logic)
  └── adapters/            # HTTP handlers, DB, external services
pkg/                        # Shared utilities
docs/                       # OpenAPI, diagramas
.github/workflows/          # CI/CD
```

## 🔍 Qué Hay Dentro

### Ejemplos Funcionales

1. **Domain Layer** (`internal/domain/user.go`)
   - Entity: `User`
   - Interface: `UserRepository`
   - Errores: `ErrUserNotFound`, etc.

2. **Usecases** (`internal/usecases/`)
   - `CreateUserUsecase` - crear usuario con validación
   - `GetUserUsecase` - obtener usuario
   - Tests incluidos

3. **Adapters** (`internal/adapters/`)
   - `UserHandler` - HTTP endpoints
   - `UserRepository` - In-memory storage (pruebas) → reemplaza con DB real
   - Middleware - logging, CORS, recovery

4. **Utilities** (`pkg/`)
   - Logger simple
   - (Agrega validators, helpers, etc.)

### Configuración

- `go.mod` / `go.sum` - Dependencias (mínimas)
- `Dockerfile` - Multi-stage para producción
- `.env.example` - Variables de entorno
- `.gitignore` - Excluye archivos locales

### CI/CD (Próximamente)

- `.github/workflows/test.yml` - Tests en PR
- `.github/workflows/deploy.yml` - Deploy en releases

## 🛠️ Desarrollo

### Agregar un Nuevo Endpoint

Sigue el patrón de `user`:

1. **Crea entity en `domain/`**
   ```go
   // internal/domain/product.go
   type Product struct { ... }
   type ProductRepository interface { ... }
   ```

2. **Crea usecase en `usecases/`**
   ```go
   // internal/usecases/create_product.go
   type CreateProductUsecase struct { ... }
   func (uc *CreateProductUsecase) Execute(...) { ... }
   ```

3. **Crea handler en `adapters/http/handlers/`**
   ```go
   // internal/adapters/http/handlers/product_handler.go
   type ProductHandler struct { ... }
   func (h *ProductHandler) CreateProduct(w, r) { ... }
   ```

4. **Crea repository en `adapters/persistence/`**
   ```go
   // internal/adapters/persistence/product_repository.go
   type ProductRepository struct { ... }
   func (r *ProductRepository) Save(...) { ... }
   ```

5. **Registra en `main.go`**
   ```go
   productRepo := persistence.NewProductRepository(log)
   createProductUC := usecases.NewCreateProductUsecase(productRepo)
   productHandler := handlers.NewProductHandler(createProductUC, log)
   mux.HandleFunc("POST /api/v1/products", productHandler.CreateProduct)
   ```

6. **Agrega tests**
   ```go
   // internal/usecases/create_product_test.go
   func TestCreateProduct(t *testing.T) { ... }
   ```

### Tests

```bash
# Todos
go test ./...

# Con cobertura
go test ./... -cover

# Específico
go test -run TestCreateUser ./internal/usecases
```

### Docker

```bash
# Build
docker build -t my-api:latest .

# Run
docker run -p 8080:8080 -e LOG_LEVEL=info my-api:latest
```

## 🔄 Git Workflow

1. Crea rama: `feature/new-feature`
2. Commit: `git commit -m "feat: add new feature"`
3. Push y abre PR a `staging`
4. Tests corren automáticamente
5. Merge después de aprobación

Ver [BRANCHING-STRATEGY.md](../epa-standards/docs/BRANCHING-STRATEGY.md) para detalles.

## 📝 Próximos Pasos

1. ✅ Lee `CLAUDE.md`
2. ✅ Run `go run ./cmd/api`
3. ✅ Revisa `internal/domain/user.go` (entity de ejemplo)
4. ✅ Revisa `internal/usecases/create_user.go` (usecase de ejemplo)
5. ✅ Revisa `internal/adapters/http/handlers/user_handler.go` (handler de ejemplo)
6. ✅ Crea tu primer endpoint copiando el patrón de `user`
7. ✅ Escribe tests
8. ✅ Abre PR a `staging`

## 📞 Preguntas?

- Lee `CLAUDE.md` (contiene respuestas a Q&A comunes)
- Revisa archivos de ejemplo en `internal/`
- Pregunta en Slack/equipo

## 📄 Estado

- ✅ Architecture: Hexagonal
- ✅ Examples: Users (CRUD)
- ✅ Tests: Included
- ✅ Docker: Configured
- ⏳ Database: TODO (reemplaza in-memory repo)
- ⏳ OpenAPI: TODO (instala swag)
- ⏳ GitHub Workflows: TODO (configura en tu repo)

---

**Última actualización:** 2026-05-20

**¡Buena suerte con tu API!** 🚀
