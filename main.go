package main

import (
	"context"
	"dubmer-bono/app/api"
	"log"
)

func main() {
	ctx := context.Background()
	logger := getLogger()

	server, err := api.NewServer(ctx, logger, "/tmp")
	if err != nil {
		log.Fatal(err)
		return
	}
	server.StartServer(ctx)
}
