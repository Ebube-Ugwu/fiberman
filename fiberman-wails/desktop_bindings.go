//go:build bindings

package main

import "net/http"

func newDesktopRuntime() (http.Handler, *App, error) {
	return http.NotFoundHandler(), NewApp(nil), nil
}
