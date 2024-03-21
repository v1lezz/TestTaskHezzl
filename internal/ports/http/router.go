package http

import (
	"context"
	"net/http"
	"testTaskHezzl/internal/app"
)

func AppRouter(ctx context.Context, a *app.App) {
	http.HandleFunc("/good/create", nil)
	http.HandleFunc("/good/update", nil)
	http.HandleFunc("/good/remove", nil)
	http.HandleFunc("/goods/list", nil)
	http.HandleFunc("/goods/reprioritiize", nil)
}
