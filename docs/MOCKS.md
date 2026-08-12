# Mocks con Mockery

Este repo usa [mockery](https://vektra.github.io/mockery/latest/) v2 para
generar mocks de las interfaces en `internal/pkg/ports`. **Nunca se editan
los archivos de `mocks/` a mano** — siempre se regeneran.

## Instalar mockery

```bash
go install github.com/vektra/mockery/v2@v2.53.3
```

Verifica la instalación:

```bash
mockery --version
```

## Configuración: `.mockery.yaml`

```yaml
with-expecter: false
mockname: "{{.InterfaceName}}"
filename: "{{.InterfaceName}}.go"
outpkg: mocks
dir: mocks

packages:
  github.com/epa-datos/epa-standards-backend/internal/pkg/ports:
    interfaces:
      ExampleRepository:
      ExampleService:
```

- `packages` lista, por paquete, **qué interfaces** generar (así no se
  mockea nada por accidente).
- `outpkg: mocks` + `dir: mocks` → todos los mocks quedan en la carpeta
  `mocks/` en la raíz del repo, como en `admin-tool-api`.
- `mockname`/`filename` con `{{.InterfaceName}}` → un archivo por interfaz,
  con el mismo nombre (`ExampleRepository.go`, `ExampleService.go`, ...),
  igual que los mocks existentes de `admin-tool-api`.

## Generar / regenerar mocks

```bash
# Desde la raíz del repo
mockery
# o
make mocks
```

Esto (re)genera `mocks/ExampleRepository.go` y `mocks/ExampleService.go` a
partir de las interfaces actuales en `internal/pkg/ports/example.go`.

## Agregar mocks para un recurso nuevo

1. Define las interfaces en `internal/pkg/ports/mi_recurso.go`.
2. Agrégalas a `.mockery.yaml`:

   ```yaml
   packages:
     github.com/epa-datos/epa-standards-backend/internal/pkg/ports:
       interfaces:
         ExampleRepository:
         ExampleService:
         MiRecursoRepository:   # 👈 nuevo
         MiRecursoService:      # 👈 nuevo
   ```

3. Corre `make mocks` (o `mockery`).
4. Verifica que aparecieron `mocks/MiRecursoRepository.go` y
   `mocks/MiRecursoService.go`.
5. Corre `go mod tidy` si es la primera vez que usas mocks en el módulo.

## Usar un mock en un test

Cada mock generado trae un constructor `mocks.NewXxx(t)` que registra
`t.Cleanup(mock.AssertExpectations)` automáticamente — si programas una
expectativa con `.On(...)` y el código no la llama, el test falla solo.

```go
func TestService_GetByID(t *testing.T) {
    repo := mocks.NewExampleRepository(t)
    svc := example.NewService(repo)

    repo.On("GetByID", mock.Anything, "1").
        Return(&entity.Example{ID: "1", Name: "Sample"}, nil)

    got, err := svc.GetByID(context.Background(), "1")

    assert.NoError(t, err)
    assert.Equal(t, "Sample", got.Name)
}
```

Ver `internal/pkg/service/example/service_test.go` y
`internal/infrastructure/api/example/handlers_test.go` para más ejemplos
(incluyendo tabla de casos y mocks de servicio para probar handlers HTTP).

## Errores comunes

**`mockery: command not found`**
→ Asegúrate que `$(go env GOPATH)/bin` esté en tu `PATH`.

**Los mocks no compilan después de cambiar una interfaz**
→ Volviste a correr `make mocks` después de editar `ports/*.go`? Los mocks
  quedan desincronizados si no se regeneran.

**`go mod tidy` falla buscando el paquete `mocks` como módulo externo**
→ Pasa si `mocks/` está vacío (sin `.go` files) pero algún `_test.go` ya lo
  importa. Corre `make mocks` **antes** de `go mod tidy`.
