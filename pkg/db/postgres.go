package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
	"pizzeria/internal/config"
)

func New(ctx context.Context, config config.Config) (*pgxpool.Pool, error) {
	pgxConfig, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=require",
		config.Db.User,
		config.Db.Pass,
		config.Db.Host,
		config.Db.Port,
		config.Db.Name,
	))
	if err != nil {
		return nil, err
	}

	pgxConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	return pgxpool.NewWithConfig(ctx, pgxConfig)
}
