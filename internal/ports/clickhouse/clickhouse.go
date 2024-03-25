package clickhouse_port

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"testTaskHezzl/internal/logger"
)

const (
	addLog = "INSERT INTO db_log.LogInfo VALUES ($1,$2,$3,$4,$5,$6,$7)"
)

type ClickhouseConn struct {
	conn driver.Conn
}

func NewClickhouseConn(ctx context.Context, port int, migrationsPath string) (*ClickhouseConn, error) {
	//if err := migrations_clickhouse.Up(port, migrationsPath); err != nil {
	//	return nil, err
	//}
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("ch:%d", port)},
		Auth: clickhouse.Auth{
			Username: "v1lezz",
			Password: "1234",
			Database: "default",
		},
	})
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(ctx); err != nil {
		return nil, err
	}
	return &ClickhouseConn{
		conn: conn,
	}, nil
}

func (cc *ClickhouseConn) AddLog(ctx context.Context, info logger.LogInfo) error {
	return cc.conn.Exec(ctx, addLog, info.ID, info.ProjectID, info.Name, info.Description, info.Priority, info.Removed, info.EventTime)
}
