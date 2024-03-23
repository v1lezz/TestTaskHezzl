package app

import (
	"context"
	"log"
	"net/http"
	"testTaskHezzl/internal/good"
)

type App struct {
	HTTPServer *http.Server
	DBRepo     good.DBRepository
	CacheRepo  good.CacheRepository
}

func NewApp(ctx context.Context) (*App, error) {
	return &App{}, nil
}

func (a *App) Run() error {
	defer a.Close()
	if err := a.HTTPServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (a *App) Close() {

}
