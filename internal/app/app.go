package app

import "context"

type App struct {
}

func NewApp(ctx context.Context) (*App, error) {
	return &App{}, nil
}

func (a *App) Run() error {
	defer a.Close()
	
}

func (a *App) Close() {

}
