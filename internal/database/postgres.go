package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kevalsabhani/keeper/internal/configs"
)

var db *pgxpool.Pool

func New(config *configs.Config) (*pgxpool.Pool, error) {
	if db != nil {
		return db, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, config.DbUrl)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}
