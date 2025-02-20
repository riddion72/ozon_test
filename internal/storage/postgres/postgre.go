package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/riddion72/ozon_test/internal/config"
	"github.com/riddion72/ozon_test/internal/logger"
)

const (
	maxRetries = 5
	retryDelay = 5 * time.Second
)

func ConnectWithRetries(ctx context.Context, cfg config.DB) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Name, cfg.Password)

	for i := 0; i < maxRetries; i++ {
		db, err = sqlx.ConnectContext(ctx, "postgres", dsn)
		if err == nil {
			if err = db.Ping(); err == nil {
				logger.Info("Successfully connected to PostgreSQL",
					slog.String("host", cfg.Host),
					slog.String("db", cfg.Name))
				return db, nil
			}
		}

		logger.Warn("Failed to connect to PostgreSQL, retrying...",
			slog.Int("attempt", i+1),
			slog.String("error", err.Error()))

		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("failed to connect after %d attempts: %v", maxRetries, err)
}
