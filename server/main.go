package main

import (
	"context"
	"dubmer-bono/app/api/udp"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	logger := getLogger()
	err := udp.NewUDPServer(ctx, logger, 4345)
	if err != nil {
		logger.Errorf("Error Setting UDP: %w", err)
		return
	}
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
