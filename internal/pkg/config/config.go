// Package config centralizes application configuration.
//
// It uses Viper to read configuration from (in order of precedence):
//  1. Real OS environment variables (e.g. those injected by Cloud Run, Docker, k8s, CI)
//  2. A local ".env" file (see .env.example) — handy for local development
//
// Real env vars always win over the file, because viper.AutomaticEnv() is
// combined with a config file that is treated as optional. This means the
// exact same binary works locally (with a .env file) and in production
// (with real secrets injected by the platform, no file involved).
package config

import (
	"log"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var (
	once sync.Once
	// Cfg holds the loaded configuration. It is populated once on first
	// access (see Load) and is safe to read concurrently afterwards.
	Cfg Config
)

// Config stores every configuration value used by the application.
// Field tags map 1:1 to environment variable names — add new fields here
// whenever you need a new env var and Viper will pick it up automatically.
type Config struct {
	// HTTP server
	ServerPort     string `mapstructure:"SERVER_PORT"`
	Environment    string `mapstructure:"ENVIRONMENT"`
	LogLevel       string `mapstructure:"LOG_LEVEL"`
	AllowedOrigins string `mapstructure:"ALLOWED_ORIGINS"`

	// Postgres (see internal/infrastructure/repositories/postgres)
	PostgresUser string `mapstructure:"POSTGRESDB_USER"`
	PostgresPass string `mapstructure:"POSTGRESDB_PASS"`
	PostgresName string `mapstructure:"POSTGRESDB_NAME"`
	PostgresHost string `mapstructure:"POSTGRESDB_HOST"`
	PostgresPort string `mapstructure:"POSTGRESDB_PORT"`

	// Firestore (see internal/infrastructure/repositories/firestore)
	FirestoreProjectID string `mapstructure:"FIRESTORE_PROJECT_ID"`

	// Auth middleware placeholder (internal/infrastructure/api/middlewares/auth.go)
	APIAuthToken string `mapstructure:"API_AUTH_TOKEN"`
}

// Load reads configuration once and returns it. It is also called from
// init() so importing this package is enough to have Cfg populated, but
// calling Load() explicitly is preferred in tests/main for clarity.
func Load() Config {
	once.Do(func() {
		viper.SetConfigName(".env")
		viper.SetConfigType("env")

		// Look for the .env file from common call sites: repo root (main.go,
		// go test ./...) and nested package directories (handlers_test.go,
		// service_test.go living a few levels deep under internal/...).
		viper.AddConfigPath(".")
		viper.AddConfigPath("../../../../")
		viper.AddConfigPath("../../../")
		viper.AddConfigPath("../../")
		viper.AddConfigPath("../")

		// Real environment variables always override file values.
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if err := viper.ReadInConfig(); err != nil {
			// The .env file is OPTIONAL: in production the platform injects
			// real env vars directly and no file exists, so we only log.
			log.Printf("config: no .env file found (%v), relying on real environment variables", err)
		}

		if err := viper.Unmarshal(&Cfg); err != nil {
			log.Fatal("config: failed to unmarshal configuration: " + err.Error())
		}

		if Cfg.ServerPort == "" {
			Cfg.ServerPort = "8080"
		}
	})
	return Cfg
}

func init() {
	Load()
}
