# Estructura de carpetas y archivos

Este repo sigue la misma organización que
[`admin-tool-api`](https://github.com/epa-datos/admin-tool-api): separa
**infraestructura** (HTTP, bases de datos) de **lógica de negocio** (`pkg`).
No es arquitectura hexagonal "de libro" con `domain/usecases/adapters`; es el
patrón que ya usa el equipo en producción, así que si conoces `admin-tool-api`
te vas a mover por este repo sin fricción.

```
.
├── main.go                          # Entry point: solo llama a api.RunServer()
├── go.mod / go.sum                  # Dependencias del módulo
├── Dockerfile                       # Build multi-stage (alpine)
├── Makefile                        # Atajos: run, test, lint, mocks, docker-*
├── .golangci.yml                    # Configuración de golangci-lint
├── .mockery.yaml                    # Qué interfaces mockear y dónde escribirlas
├── .env.example                     # Plantilla de variables de entorno (sin secretos)
├── .gitignore
├── .github/workflows/                # CI: run_tests.yml, golangci-lint.yml, cloudrun_deploy.yml
├── docs/                             # Estás aquí
├── mocks/                            # Generado por mockery — nunca se edita a mano
└── internal/
    ├── infrastructure/               # 🔵 Todo lo que habla con el mundo exterior
    │   ├── api/                      # Capa HTTP (gin)
    │   │   ├── server.go             #   gin.Engine, middlewares globales, RunServer()
    │   │   ├── routes.go             #   Dependency injection + montaje de cada recurso
    │   │   ├── middlewares/          #   auth.go, (agrega cors, logging, etc. aquí)
    │   │   └── example/              #   👉 Vertical slice de ejemplo: handlers.go, routes.go, handlers_test.go
    │   └── repositories/              # Implementaciones de los ports.*Repository
    │       ├── postgres/              #   👉 client.go (conexión GORM) + example_repository.go
    │       └── firestore/             #   👉 client.go (conexión Firestore) + example_repository.go
    └── pkg/                          # 🟡 Lógica de negocio — no debería importar `gin`, `gorm`, etc.
        ├── config/                   #   Viper: carga .env / variables de entorno reales
        │   └── config.go
        ├── entity/                   #   Structs de dominio (lo que viaja en JSON / se persiste)
        │   └── example.go
        ├── ports/                    #   Interfaces: el contrato entre service y repository/handler
        │   └── example.go
        ├── service/                  #   Un paquete por recurso, con su lógica + tests
        │   └── example/
        │       ├── service.go
        │       └── service_test.go
        └── utils/                    #   Helpers puros y compartidos (paginación, fechas, ...)
            └── pagination.go
```

## Regla de dependencias

```
internal/infrastructure/api/<recurso>   →  internal/pkg/ports (interfaces)
internal/infrastructure/repositories/*  →  internal/pkg/ports (implementa la interfaz)
internal/pkg/service/<recurso>          →  internal/pkg/ports + internal/pkg/entity
internal/pkg/entity                     →  nada (solo stdlib)
```

- **`pkg/entity`** no importa nada del resto del proyecto: son datos puros.
- **`pkg/ports`** define interfaces (`ExampleRepository`, `ExampleService`).
  Es la única pieza que conocen tanto la capa HTTP como los repositorios.
- **`pkg/service/<recurso>`** implementa la interfaz de servicio y depende
  **solo de la interfaz de repositorio** (`ports.ExampleRepository`), nunca de
  `gorm.DB` o `firestore.Client` directamente. Esto es lo que permite:
  - cambiar Postgres ↔ Firestore sin tocar el servicio (ver `routes.go`), y
  - probar el servicio con un mock en vez de una base de datos real.
- **`infrastructure/repositories/<engine>`** implementa `ports.ExampleRepository`
  usando una tecnología concreta.
- **`infrastructure/api/<recurso>`** traduce HTTP ↔ `ports.ExampleService`.

## ¿Por qué un recurso "Example" en vez de vacío?

Un repo de ejemplo vacío obliga a adivinar el patrón. Con `Example`
implementado de punta a punta (entity → ports → service con tests → dos
repositorios → handlers con tests → rutas) puedes **copiar la carpeta
completa**, renombrar `Example` → `TuRecurso` y ya tienes las cinco piezas
en el lugar correcto. Bórralo cuando ya no lo necesites como referencia.

## Dónde tocar cuando agregas un recurso

Ver la sección *"Cómo agregar un recurso nuevo"* en [CLAUDE.md](../CLAUDE.md)
y en el [README](../README.md#-cómo-agregar-un-recurso-nuevo).
