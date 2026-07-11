package api

import (
	"context"
	"dubmer-bono/app/service"
	ingestion "dubmer-bono/app/service/ingestionservice"
	monitor "dubmer-bono/app/service/monitorservice"
	"dubmer-bono/app/service/radio"
	"dubmer-bono/app/types"
	"dubmer-bono/app/types/entity"

	"github.com/gordonklaus/portaudio"
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

	msg_chan := make(chan entity.RadioMessage) //To send messages to radio from MonitorService

	monitor, err := monitor.NewService(ctx, logger, path, repo, msg_chan)
	if err != nil {
		return Services{}, err
	}

	err = portaudio.Initialize() //So Radio can use Speaker and Mic
	if err != nil {
		return Services{}, err
	}
	radio, err := radio.NewService(ctx, logger, path, repo, monitor.GetPressure, msg_chan)
	if err != nil {
		return Services{}, err
	}
	return Services{ingestion: ingestion, monitor: monitor, radio: radio}, nil
}
