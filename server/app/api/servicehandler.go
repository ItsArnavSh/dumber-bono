package api

import (
	"context"
	ingestion "dubmer-bono/app/service/ingestionservice"
	"dubmer-bono/app/types"

	"go.uber.org/zap"
)

type Services struct {
	ingestion types.Ingestion
	monitor   types.Monitor
	radio     types.Radio
	agent     types.Agent
}

func InitServices(ctx context.Context, path string, logger *zap.SugaredLogger) (Services, error) {
	ingestion, err := ingestion.NewService(ctx, logger, path)
	if err != nil {
		return Services{}, err
	}
	return Services{ingestion: ingestion}, nil
}
