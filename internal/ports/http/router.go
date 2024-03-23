package http

import (
	"context"
	"net/http"
	"testTaskHezzl/internal/app"
)

func AppRouter(ctx context.Context, a *app.App) {
	http.HandleFunc("/good/create", CreateGoodHandler(ctx, a))
	http.HandleFunc("/good/update", UpdateGoodHandler(ctx, a))
	http.HandleFunc("/good/remove", RemoveGoodHandler(ctx, a))
	http.HandleFunc("/goods/list", GetListGoodsHandler(ctx, a))
	http.HandleFunc("/goods/reprioritiize", ReprioritiizeGoodHandler(ctx, a))
}
