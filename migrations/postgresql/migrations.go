package migrations_postgresql

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
)

func Up(pool *pgxpool.Pool, path string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	db := stdlib.OpenDBFromPool(pool)
	if err := goose.Up(db, path); err != nil {
		return err
	}
	if err := db.Close(); err != nil {
		return err
	}
	return nil
}
