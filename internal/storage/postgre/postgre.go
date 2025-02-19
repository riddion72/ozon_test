package storage

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/riddion72/ozon_test/internal/config"
)

const (
	col1 string = "article_name"
	col2 string = "article_content"
)

type pgRepo struct {
	db *sqlx.DB
}

func New(conn *sqlx.DB) *pgRepo {
	return &pgRepo{db: conn}
}

func NewConnection(ctx context.Context, cfg config.DB) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Name, cfg.Password)

	db, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
