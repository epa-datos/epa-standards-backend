# Variables de entorno y configuración (Viper)

## Cómo funciona

`internal/pkg/config/config.go` usa [Viper](https://github.com/spf13/viper)
para poblar un struct `Config` a partir de:

1. **Variables de entorno reales del sistema** (`viper.AutomaticEnv()`) —
   siempre tienen prioridad. Así funciona en Cloud Run/Docker/k8s, donde el
   valor real llega inyectado por la plataforma y **no existe ningún archivo**.
2. **Un archivo `.env` local** (opcional) — conveniente para desarrollo. Si no
   existe, solo se registra un log informativo y el programa sigue: nunca
   se cae por falta de `.env` en producción.

```go
func Load() Config {
    once.Do(func() {
        viper.SetConfigName(".env")
        viper.SetConfigType("env")
        viper.AddConfigPath(".")
        // ...más rutas para que funcione también desde tests anidados

        viper.AutomaticEnv() // las env vars reales siempre ganan

        if err := viper.ReadInConfig(); err != nil {
            log.Printf("config: no .env file found (%v), relying on real environment variables", err)
        }
        if err := viper.Unmarshal(&Cfg); err != nil {
            log.Fatal("config: failed to unmarshal configuration: " + err.Error())
        }
    })
    return Cfg
}
```

> **Diferencia con `admin-tool-api`:** ese repo usa `log.Fatal` si no
> encuentra el archivo `app.env`, lo cual obliga a tener siempre un archivo
> presente (y en ese repo terminó commiteado con secretos reales — no lo
> repitas). Aquí el archivo es genuinamente opcional para poder correr en
> contenedores donde las variables llegan del entorno real.

## Agregar una variable nueva

1. Agrega el campo a `Config` en `internal/pkg/config/config.go` con su tag
   `mapstructure`:

   ```go
   type Config struct {
       // ...
       MyNewSetting string `mapstructure:"MY_NEW_SETTING"`
   }
   ```

2. Documéntala en `.env.example` con un valor de ejemplo (nunca un secreto real):

   ```env
   MY_NEW_SETTING=some-placeholder
   ```

3. Úsala donde la necesites vía `config.Cfg.MyNewSetting` (el paquete la
   carga automáticamente en su `init()`, así que solo con importar
   `internal/pkg/config` ya está disponible).

## Variables actuales

| Variable | Descripción | Dónde se usa |
|---|---|---|
| `SERVER_PORT` | Puerto HTTP (default `8080`) | `api/server.go` |
| `ENVIRONMENT` | `development` / `production` — cambia el modo de gin | `api/server.go` |
| `LOG_LEVEL` | Nivel de log (placeholder para tu logger) | — |
| `ALLOWED_ORIGINS` | Orígenes permitidos por CORS, separados por coma | `api/server.go` |
| `POSTGRESDB_USER/PASS/NAME/HOST/PORT` | Conexión a Postgres | `repositories/postgres/db.go` |
| `FIRESTORE_PROJECT_ID` | Proyecto GCP para el cliente de Firestore | `repositories/firestore/client.go` |
| `API_AUTH_TOKEN` | Token estático que exige el middleware de auth (vacío = sin auth, solo para dev) | `api/middlewares/auth.go` |

## Reglas de seguridad

- `.env` está en `.gitignore`. **Solo `.env.example` se commitea**, y solo
  con placeholders.
- Nunca pongas tokens, contraseñas o JSON de service account reales en
  `.env.example`, en el código, ni en los workflows de `.github/`.
- En CI/CD, los valores reales se inyectan como *secrets* de GitHub Actions
  o *secrets* de Google Secret Manager (ver el `TODO` en
  `.github/workflows/cloudrun_deploy.yml`), nunca como archivo commiteado.
- Si necesitas credenciales de Google Cloud localmente, usa
  `gcloud auth application-default login` en vez de descargar un JSON de
  service account al repo.
