package main

import (
	"context"
	"log"
	"testTaskHezzl/internal/app"
	"testTaskHezzl/internal/config"
	clickhouse_port "testTaskHezzl/internal/ports/clickhouse"
	"testTaskHezzl/internal/ports/http"
	nats_port "testTaskHezzl/internal/ports/nats"
	"testTaskHezzl/internal/ports/postgres"
	"testTaskHezzl/internal/ports/redis"
)

func main() {
	ctx := context.Background()
	dbCfg := config.DBConfig{
		User:     "v1lezz",
		Password: "1234",
		Database: "goods",
		Host:     "pg_db",
		Port:     5432,
	}
	httpSrv := http.NewHTTPServer(8080)
	pg, err := postgres.NewConnection(ctx, dbCfg, "/goods/migrations/postgresql")
	if err != nil {
		log.Fatal(err)
	}
	redis := redis.NewRedisConn(ctx)
	ch, err := clickhouse_port.NewClickhouseConn(ctx, 9000, "/goods/migrations/clickhouse")
	if err != nil {
		log.Fatal("ch run: ", err)
	}
	logger, err := nats_port.NewNatsConn(ctx, ch)
	if err != nil {
		log.Fatal("nats run: ", err)
	}

	a, err := app.NewApp(httpSrv, pg, redis, logger)
	if err != nil {
		log.Fatal("app new run: ", err)
	}
	http.AppRouter(ctx, a)
	if err = a.Run(); err != nil {
		log.Fatal("server run: ", err)
	}
}
