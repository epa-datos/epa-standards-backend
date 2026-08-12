# EPA Digital — Standard Backend API (Go)

Plantilla estándar de API en Go para EPA Digital. Está construida con la
**misma estructura de carpetas que [admin-tool-api](https://github.com/epa-datos/admin-tool-api)**,
pero con la lógica de negocio reemplazada por un recurso de ejemplo genérico
(`Example`) para que sirva como punto de partida limpio.

**Este repo está pensado para hacer `fork` / `copy` y empezar un proyecto nuevo.**
No contiene lógica de negocio real de ningún cliente ni credenciales — todo lo
sensible se reemplazó por placeholders.

## 🚀 Quick start

```bash
# 1. Clona (o haz fork) y renombra el módulo si vas a publicarlo con otro nombre
git clone https://github.com/epa-datos/epa-standards-backend.git mi-nueva-api
cd mi-nueva-api

# 2. Copia las variables de entorno
cp .env.example .env

# 3. Descarga dependencias
go mod download

# 4. Corre el servidor
go run .
# Deberías ver: "Server running on :8080"

# 5. Prueba el endpoint de ejemplo
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/examples
```

> El endpoint de ejemplo usa Postgres por defecto (`internal/infrastructure/api/routes.go`).
> Si no tienes Postgres corriendo localmente, el servidor arranca pero las
> rutas de `/api/v1/examples` fallarán al conectar — es esperado, es solo un
> ejemplo. Ver [docs/VARIABLES-ENTORNO.md](docs/VARIABLES-ENTORNO.md).

## 📚 Documentación

| Documento | Contenido |
|---|---|
| [CLAUDE.md](./CLAUDE.md) | Guía completa de arquitectura, convenciones y cómo agregar un recurso nuevo |
| [docs/ESTRUCTURA.md](docs/ESTRUCTURA.md) | Qué hay en cada carpeta y por qué |
| [docs/MOCKS.md](docs/MOCKS.md) | Cómo generar/regenerar mocks con **mockery** |
| [docs/TESTING.md](docs/TESTING.md) | Cómo escribir y correr pruebas unitarias |
| [docs/VARIABLES-ENTORNO.md](docs/VARIABLES-ENTORNO.md) | Variables de entorno, Viper, `.env` vs entorno real |

## 📁 Estructura (resumen)

```
├── main.go                                  # Entry point
├── go.mod / go.sum
├── Dockerfile
├── Makefile                                 # run, test, lint, mocks, docker-*
├── .golangci.yml
├── .mockery.yaml                            # Config de mockery
├── .env.example
├── .github/workflows/                       # CI: tests, lint, deploy
├── docs/                                    # Documentación detallada
├── mocks/                                   # Mocks generados por mockery (NO editar a mano)
└── internal/
    ├── infrastructure/
    │   ├── api/
    │   │   ├── server.go                    # gin.Engine, middlewares globales, Run()
    │   │   ├── routes.go                    # DI + montaje de cada recurso
    │   │   ├── middlewares/                 # auth, cors, etc.
    │   │   └── example/                     # 👉 Vertical slice de ejemplo (handlers + routes + tests)
    │   └── repositories/
    │       ├── postgres/                    # 👉 Implementación de ExampleRepository con GORM
    │       └── firestore/                   # 👉 Implementación de ExampleRepository con Firestore
    └── pkg/
        ├── config/                          # Viper: carga variables de entorno
        ├── entity/                          # Structs de dominio (Example, ...)
        ├── ports/                           # Interfaces (contratos) — lo que se mockea
        ├── service/example/                 # Lógica de negocio + tests
        └── utils/                           # Helpers compartidos (paginación, ...)
```

Ver el detalle completo en [docs/ESTRUCTURA.md](docs/ESTRUCTURA.md).

## 🧱 Cómo agregar un recurso nuevo

El patrón completo está implementado para `example` — cópialo:

1. **Entity** → `internal/pkg/entity/mi_recurso.go`
2. **Ports** (interfaces repo + service) → `internal/pkg/ports/mi_recurso.go`
3. **Service** (lógica de negocio) → `internal/pkg/service/mi_recurso/service.go` + `service_test.go`
4. **Repository** (Postgres o Firestore) → `internal/infrastructure/repositories/<engine>/mi_recurso_repository.go`
5. **Handlers + routes** → `internal/infrastructure/api/mi_recurso/{handlers,routes}.go` + `handlers_test.go`
6. **Wire todo** en `internal/infrastructure/api/routes.go`
7. **Genera los mocks**: agrega las interfaces nuevas a `.mockery.yaml` y corre `make mocks`
8. **Corre los tests**: `make test`

Ver el paso a paso con código en [CLAUDE.md](./CLAUDE.md#-cómo-agregar-un-recurso-nuevo).

## 🧪 Comandos comunes

```bash
make run             # go run .
make test            # go test ./... -v
make test-cover      # go test ./... -cover
make lint            # golangci-lint run
make mocks           # regenera mocks/ con mockery
make build           # compila el binario
make docker-build    # docker build
make docker-run      # docker run --env-file .env
```

## 🔐 Seguridad

- **Nunca** commitees `.env` (está en `.gitignore`) ni archivos de credenciales
  (`*-service-account.json`, `*.pem`, `*.key`).
- `.env.example` solo contiene placeholders, nunca valores reales.
- El middleware `internal/infrastructure/api/middlewares/auth.go` es un
  placeholder mínimo — reemplázalo por tu proveedor real (Firebase, Auth0,
  JWT, API Gateway, ...) antes de ir a producción.

## 🔄 Git workflow

Mismo flujo que el resto de repos de EPA Digital:

```
feature/nombre-feature
fix/nombre-del-bug
refactor/algo
```

1. Crea rama desde `staging`
2. Commits: `feat: ...`, `fix: ...`, `test: ...`, `docs: ...`
3. Abre PR a `staging`, espera CI (tests + lint) y 1 aprobación
4. Merge

## 📄 Estado de este template

- ✅ Estructura idéntica a `admin-tool-api` (hexagonal-ish: api / repositories / pkg)
- ✅ Recurso de ejemplo end-to-end (`Example`): entity → ports → service → repo (Postgres y Firestore) → handlers → routes
- ✅ Mocks generados con `mockery` (`.mockery.yaml` + `mocks/`)
- ✅ Pruebas unitarias de servicio y de handler usando los mocks
- ✅ Viper + variables de entorno (`.env.example`, `internal/pkg/config`)
- ✅ Docker multi-stage + GitHub Actions (tests, lint, deploy de referencia)
- ⏳ Reemplaza `Example` por tu dominio real
- ⏳ Reemplaza el middleware de auth por uno real
- ⏳ Completa los `TODO` de `.github/workflows/cloudrun_deploy.yml` con tu proyecto GCP

---

**Última actualización:** 2026-08-11
