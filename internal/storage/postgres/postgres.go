package postgres

import (
	"context"
	"ls-0/arti/order/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, cfg *config.Config) *PostgresStorage {

	dbUrl := cfg.Pg.Url

	dbConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		return nil
	}

	// Customize pool config
	dbConfig.MaxConns = cfg.Pg.Pool.MaxConns
	dbConfig.MinConns = cfg.Pg.Pool.MinConns
	dbConfig.MaxConnLifetime = cfg.Pg.Pool.MaxConnLifetime
	dbConfig.HealthCheckPeriod = cfg.Pg.Pool.HealthcheckPeriod
	dbConfig.MinIdleConns = cfg.Pg.Pool.MinIdle

	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil
	}

	return &PostgresStorage{pool: pool}
}
