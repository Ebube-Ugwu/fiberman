package main

import (
	"log"
	"net/http"

	fiberman "github.com/fiberman/fiberman-go-backend"
)

func main() {
	config := fiberman.LoadConfigFromEnv()
	server := fiberman.NewServer(config)

	log.Printf("fiberman-go-backend listening on :%s", config.ServerPort)
	if err := http.ListenAndServe(":"+config.ServerPort, server.Routes()); err != nil {
		log.Fatal(err)
	}
}
