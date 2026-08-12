# CLAUDE.md — EPA Digital Standard Backend API

Este es el template estándar de API Go de EPA Digital. Tiene la **misma
estructura de carpetas que `admin-tool-api`**, con la lógica de negocio real
reemplazada por un recurso de ejemplo (`Example`) de punta a punta. Léelo
completo antes de empezar a desarrollar.

## 🚀 Comienza aquí (5 min)

```bash
# 1. Copia variables de entorno
cp .env.example .env

# 2. Descarga dependencias
go mod download

# 3. Inicia el servidor
go run .
# Deberías ver: "Server running on :8080"

# 4. Prueba
curl http://localhost:8080/health
```

---

## 📋 Comandos comunes

```bash
# Desarrollo
go run .                      # Inicia servidor localmente
make run                      # Igual, vía Makefile

# Testing
go test ./...                 # Todos los tests
go test ./... -v              # Tests verbosos
go test ./... -cover          # Con cobertura
go test -run TestNombreTest ./internal/pkg/service/example
make test / make test-cover

# Mocks (ver docs/MOCKS.md)
mockery                       # Regenera mocks/ desde internal/pkg/ports
make mocks

# Linting
go vet ./...
golangci-lint run
make lint

# Building
go build -o epa-standards-backend .
make build

# Docker
make docker-build
make docker-run
```

---

## 🏗️ Arquitectura

**Principio:** la lógica de negocio (`internal/pkg`) es independiente de
HTTP y de la base de datos concreta. Ambos lados solo se conocen a través de
interfaces (`internal/pkg/ports`).

```
┌──────────────────────────────────────────────┐
│  internal/infrastructure/api/<recurso>        │  ← Handlers gin (HTTP)
├──────────────────────────────────────────────┤
│  internal/pkg/service/<recurso>               │  ← Lógica de negocio
├──────────────────────────────────────────────┤
│  internal/pkg/ports                           │  ← Interfaces (el contrato)
├──────────────────────────────────────────────┤
│  internal/infrastructure/repositories/<engine> │  ← Postgres / Firestore
└──────────────────────────────────────────────┘
```

**Flujo real: crear un Example**

```
1. HTTP POST /api/v1/examples {"name": "..."}
2. example.Handler.Create()              ← capa HTTP (internal/infrastructure/api/example)
3. exampleSvc.service.Create()           ← lógica de negocio (internal/pkg/service/example)
4. Valida el nombre (regla de negocio)
5. postgres.ExampleRepository.Create()   ← implementación concreta (inyectada por interfaz)
6. Respuesta 201 con el Example creado
```

### Estructura de directorios

Ver el detalle completo, con la explicación de cada carpeta, en
[docs/ESTRUCTURA.md](docs/ESTRUCTURA.md). Resumen:

```
internal/
├── infrastructure/
│   ├── api/
│   │   ├── server.go        # gin.Engine + Run()
│   │   ├── routes.go        # DI: repo → service → handler → rutas
│   │   ├── middlewares/     # auth.go (placeholder — reemplázalo)
│   │   └── example/         # handlers.go, routes.go, handlers_test.go
│   └── repositories/
│       ├── postgres/        # client.go (conexión GORM) + example_repository.go
│       └── firestore/       # client.go (conexión Firestore) + example_repository.go
└── pkg/
    ├── config/               # Viper — ver docs/VARIABLES-ENTORNO.md
    ├── entity/               # Example, ExamplesResponse
    ├── ports/                # ExampleRepository, ExampleService (interfaces)
    ├── service/example/      # Lógica de negocio + service_test.go
    └── utils/                # pagination.go
mocks/                        # Generado por mockery — ver docs/MOCKS.md
```

### Reglas de oro

**✅ DO:**
- `pkg/service/<recurso>` depende de `ports.<Recurso>Repository` (interfaz),
  nunca de `gorm.DB` o `firestore.Client` directamente.
- Los handlers HTTP dependen de `ports.<Recurso>Service` (interfaz).
- Cada repositorio concreto (`postgres/`, `firestore/`) implementa la misma
  interfaz de `ports/` — son intercambiables (ver `routes.go`).
- Todo error de negocio se declara como sentinel error (`var ErrX = errors.New(...)`)
  en el paquete del servicio, y el handler lo traduce a un status HTTP en
  `handleServiceError`.

**❌ DON'T:**
- Un handler no debe llamar directo a la base de datos (Handler → DB).
- `internal/pkg/entity` no debe importar nada de `internal/infrastructure`.
- Un servicio no debe conocer detalles HTTP (status codes, headers, gin.Context).

---

## 🔌 Handlers & HTTP (gin)

### Patrón de handler

```go
// internal/infrastructure/api/example/handlers.go

type Handler struct {
    svc ports.ExampleService // interfaz, no la implementación concreta
}

func NewHandler(svc ports.ExampleService) *Handler {
    return &Handler{svc: svc}
}

func (h *Handler) Create(c *gin.Context) {
    var req entity.Example
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.svc.Create(c.Request.Context(), &req); err != nil {
        h.handleServiceError(c, err) // traduce sentinel errors → status HTTP
        return
    }

    c.JSON(http.StatusCreated, req)
}
```

### Registrar rutas

```go
// internal/infrastructure/api/example/routes.go
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
    rg.GET("", h.List)
    rg.GET("/:id", h.GetByID)
    rg.POST("", h.Create)
    rg.PUT("/:id", h.Update)
    rg.DELETE("/:id", h.Delete)
}
```

```go
// internal/infrastructure/api/routes.go — un bloque de 3 líneas por recurso
exampleRepo := postgres.NewExampleRepository(postgres.NewClient())
exampleService := exampleSvc.NewService(exampleRepo)
exampleHandler := exampleAPI.NewHandler(exampleService)
exampleAPI.RegisterRoutes(v1.Group("/examples"), exampleHandler)
```

---

## 🧱 Cómo agregar un recurso nuevo

Copia el patrón de `example` reemplazando el nombre por el de tu dominio
(`Invoice`, `Client`, `Campaign`, ...):

1. **Entity** — `internal/pkg/entity/invoice.go`
   ```go
   type Invoice struct {
       ID     string `json:"id"`
       Amount int    `json:"amount"`
   }
   ```

2. **Ports** — `internal/pkg/ports/invoice.go`
   ```go
   type InvoiceRepository interface {
       GetByID(ctx context.Context, id string) (*entity.Invoice, error)
       Create(ctx context.Context, invoice *entity.Invoice) error
       // ...
   }
   type InvoiceService interface {
       GetByID(ctx context.Context, id string) (*entity.Invoice, error)
       Create(ctx context.Context, invoice *entity.Invoice) error
       // ...
   }
   ```

3. **Service** — `internal/pkg/service/invoice/service.go` (+ `service_test.go`
   copiando `internal/pkg/service/example/service_test.go`)

4. **Repository** — `internal/infrastructure/repositories/postgres/invoice_repository.go`
   (o `firestore/`), implementando `ports.InvoiceRepository`

5. **Handlers + routes** — `internal/infrastructure/api/invoice/{handlers,routes}.go`
   (+ `handlers_test.go` copiando el de `example`)

6. **Wire** en `internal/infrastructure/api/routes.go`:
   ```go
   invoiceRepo := postgres.NewInvoiceRepository(postgres.NewClient())
   invoiceService := invoiceSvc.NewService(invoiceRepo)
   invoiceHandler := invoiceAPI.NewHandler(invoiceService)
   invoiceAPI.RegisterRoutes(v1.Group("/invoices"), invoiceHandler)
   ```

7. **Mocks** — agrega `InvoiceRepository`/`InvoiceService` a `.mockery.yaml`
   y corre `make mocks` (ver [docs/MOCKS.md](docs/MOCKS.md))

8. **Tests** — `make test`

---

## 🧪 Testing

Ver [docs/TESTING.md](docs/TESTING.md) para el detalle completo con ejemplos.
Resumen:

- **Servicio** (`pkg/service/<recurso>/service_test.go`): mockea
  `ports.<Recurso>Repository`, prueba solo lógica de negocio.
- **Handler** (`infrastructure/api/<recurso>/handlers_test.go`): mockea
  `ports.<Recurso>Service`, prueba solo la traducción HTTP ↔ servicio, con
  `httptest`.

### Coverage objetivo
- `pkg/service/*`: >80%
- `infrastructure/api/*` (handlers): >60%

---

## 🧬 Mocks (mockery)

Ver [docs/MOCKS.md](docs/MOCKS.md) para el detalle. Resumen:

```bash
go install github.com/vektra/mockery/v2@v2.53.3   # una vez
mockery                                            # regenerar mocks/
make mocks                                         # atajo
```

Las interfaces a mockear se declaran en `.mockery.yaml`. **Nunca edites a
mano** los archivos dentro de `mocks/`.

---

## 🐳 Docker & Cloud Run

```bash
# Build local
make docker-build

# Run local
make docker-run   # usa --env-file .env

# Deploy: ver .github/workflows/cloudrun_deploy.yml (tiene TODOs a completar
# con tu proyecto de GCP — está basado en el workflow real de admin-tool-api)
```

---

## 🔐 Variables de entorno

Ver [docs/VARIABLES-ENTORNO.md](docs/VARIABLES-ENTORNO.md) para el detalle
de Viper, precedencia (`.env` vs entorno real) y cómo agregar una variable
nueva. Resumen:

```bash
cp .env.example .env
```

**Nunca commitees `.env`** (está en `.gitignore`) ni pongas secretos reales
en `.env.example`.

---

## 🔄 Git Workflow

### Branch naming
```
feature/nombre-feature
fix/nombre-del-bug
refactor/algo
```

### PR process
1. Crea rama desde `staging`
2. Haz commits significativos
3. Push a remote
4. Abre PR a `staging`
5. Espera CI (tests + lint, ver `.github/workflows/`)
6. Consigue 1 aprobación
7. Merge

### Commits
```bash
git commit -m "feat: add invoice endpoint"
git commit -m "fix: validate invoice amount"
git commit -m "test: add invoice service tests"
git commit -m "docs: update API structure"
```

---

## 💡 Errores comunes & soluciones

**`module not found` / imports rotos**
```bash
go mod tidy
go mod download
```

**Tests fallan con "connection refused"**
→ Estás intentando usar una DB real en un test. Usa el mock de
  `ports.<Recurso>Repository` (ver [docs/MOCKS.md](docs/MOCKS.md)).

**`go mod tidy` falla buscando `.../mocks` como módulo externo**
→ Corre `make mocks` primero; `go mod tidy` necesita que `mocks/` ya tenga
  archivos `.go` para reconocerlo como paquete interno.

**Handler devuelve 500 en vez del status esperado**
→ Revisa `handleServiceError` en el handler: cada sentinel error nuevo del
  servicio necesita su `case errors.Is(err, ...)` ahí.

**El servidor arranca pero `/api/v1/examples` falla**
→ Es el comportamiento esperado sin Postgres corriendo — el ejemplo usa una
  base real. Levanta Postgres localmente o cambia `routes.go` para usar
  `firestore.NewExampleRepository(...)`.

**El workflow `golangci-lint` falla con `jsonschema: "linters-settings" does
not validate ...`**
→ `enable`/`disable` van bajo la sección `linters:`, y la lista de
  analizadores individuales de `go vet` va bajo `linters-settings.govet.enable`
  — nunca directo bajo `linters-settings`. Corre `golangci-lint config
  verify` localmente antes de hacer push para detectarlo (instala con
  `go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8`
  para usar la misma versión que CI).

---

## 📚 Recursos

- **Go Docs:** https://golang.org/doc/
- **Gin:** https://gin-gonic.com/docs/
- **GORM:** https://gorm.io/docs/
- **Mockery:** https://vektra.github.io/mockery/latest/
- **Viper:** https://github.com/spf13/viper
- **admin-tool-api** (repo del que se tomó esta estructura): https://github.com/epa-datos/admin-tool-api

---

## 🎯 Siguientes pasos

1. ✅ Lee esta guía completa y [docs/ESTRUCTURA.md](docs/ESTRUCTURA.md)
2. ✅ `cp .env.example .env && go run .` y verifica que corre
3. ✅ Revisa el recurso `example` de punta a punta (entity → ports → service → repo → handler)
4. ✅ Corre `go test ./... -v` y revisa cómo se usan los mocks
5. ✅ Copia el patrón de `example` para tu primer recurso real
6. ✅ Agrega sus mocks (`.mockery.yaml` + `make mocks`) y sus tests
7. ✅ Borra `example` cuando ya no lo necesites como referencia
8. ✅ Abre tu primer PR a `staging`

**¡Buena suerte!** 🚀
