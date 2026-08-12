package postgres

import (
	"fmt"
	"sync"
	"time"

	"github.com/epa-datos/epa-standards-backend/internal/pkg/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	client     *gorm.DB
	clientOnce sync.Once
)

// NewClient returns a singleton *gorm.DB connection built from
// internal/pkg/config. It retries a few times on startup because in
// containerized environments the DB may not be reachable yet on the very
// first attempt (e.g. docker-compose starting Postgres and the API
// simultaneously).
func NewClient() *gorm.DB {
	clientOnce.Do(func() {
		cfg := config.Cfg
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			cfg.PostgresHost, cfg.PostgresUser, cfg.PostgresPass, cfg.PostgresName, cfg.PostgresPort,
		)

		const maxAttempts = 10
		for attempt := 1; attempt <= maxAttempts; attempt++ {
			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err == nil {
				client = db
				logrus.Info("postgres: connected")
				return
			}

			logrus.Warnf("postgres: connection attempt %d/%d failed: %v", attempt, maxAttempts, err)
			time.Sleep(2 * time.Second)
		}

		logrus.Fatal("postgres: could not connect after ", maxAttempts, " attempts")
	})

	return client
}
