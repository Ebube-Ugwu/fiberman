package main

import (
	"context"
	"net/http"
)

type App struct {
	ctx       context.Context
	apiServer *http.Server
}

func NewApp(apiServer *http.Server) *App {
	return &App{apiServer: apiServer}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(context.Context) {
	if a.apiServer != nil {
		_ = a.apiServer.Close()
	}
}
