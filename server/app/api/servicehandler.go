package api

import (
	"context"
	"dubmer-bono/app/service"
	ingestion "dubmer-bono/app/service/ingestionservice"
	monitor "dubmer-bono/app/service/monitorservice"
	"dubmer-bono/app/types"

	"go.uber.org/zap"
)

type Services struct {
	ingestion types.Ingestion
	monitor   types.Monitor
	radio     types.Radio
	agent     types.Agent
}

func InitServices(ctx context.Context, path string, logger *zap.SugaredLogger, repo *service.Repository) (Services, error) {
	ingestion, err := ingestion.NewService(ctx, logger, path, repo)
	if err != nil {
		return Services{}, err
	}
	monitor, err := monitor.NewService(ctx, logger, path, repo)
	if err != nil {
		return Services{}, err
	}
	return Services{ingestion: ingestion, monitor: monitor}, nil
}
