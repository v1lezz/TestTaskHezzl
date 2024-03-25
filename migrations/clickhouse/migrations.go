package migrations_clickhouse

import (
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/pressly/goose"
)

func Up(port int, path string) error {
	if err := goose.SetDialect("clickhouse"); err != nil {
		return err
	}
	db := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("ch:%d", port)},
		Auth: clickhouse.Auth{
			Username: "v1lezz",
			Password: "1234",
		},
	})
	if err := goose.Up(db, path); err != nil {
		return err
	}
	if err := db.Close(); err != nil {
		return err
	}
	return nil
}
