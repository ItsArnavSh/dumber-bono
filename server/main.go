package main

import (
	"context"
	"dubmer-bono/app/api/udp"
	ingestion "dubmer-bono/app/service/ingestionservice"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	logger := getLogger()
	service, err := ingestion.NewService("/tmp")
	err = udp.NewUDPServer(ctx, logger, 4345, service)
	if err != nil {
		logger.Errorf("Error Setting UDP: %w", err)
		return
	}
	//speechservice.StartListner(ctx)
	select {}
}

func getLogger() *zap.SugaredLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	sugar := logger.Sugar()
	sugar.Info("Setting up the logger")
	return sugar
}
