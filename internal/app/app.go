package app

import (
	"log"
	"net/http"
	"testTaskHezzl/internal/good"
	"testTaskHezzl/internal/logger"
)

type App struct {
	HTTPServer *http.Server
	DBRepo     good.DBRepository
	CacheRepo  good.CacheRepository
	Logger     logger.Logger
}

func NewApp(httpSrv *http.Server, dbrepo good.DBRepository, cacherepo good.CacheRepository, logger logger.Logger) (*App, error) {
	return &App{
		HTTPServer: httpSrv,
		DBRepo:     dbrepo,
		CacheRepo:  cacherepo,
		Logger:     logger,
	}, nil
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
