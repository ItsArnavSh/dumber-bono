package main

import (
	"context"
	"dubmer-bono/app/api/udp"
	ingestion "dubmer-bono/app/service/ingestionservice"
)

func main() {
	ctx := context.Background()
	logger := getLogger()
	service, err := ingestion.NewService(ctx, "/tmp")
	err = udp.NewUDPServer(ctx, logger, 4345, service)
	if err != nil {
		logger.Errorf("Error Setting UDP: %w", err)
		return
	}
	//speechservice.StartListner(ctx)
	select {}
}
