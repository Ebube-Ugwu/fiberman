//go:build !bindings

package main

import (
	"log"
	"net"
	"net/http"

	fiberman "github.com/fiberman/fiberman-go-backend"
)

func newDesktopRuntime() (http.Handler, *App, error) {
	config := fiberman.LoadConfigFromEnv()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, nil, err
	}

	config.PlaygroundBaseURL = "http://" + listener.Addr().String()
	backend := fiberman.NewServer(config)
	apiHandler := backend.Routes()

	apiServer := &http.Server{Handler: apiHandler}
	go func() {
		if err := apiServer.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Printf("desktop api server stopped: %v", err)
		}
	}()

	return apiHandler, NewApp(apiServer), nil
}
