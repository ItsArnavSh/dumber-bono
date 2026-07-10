package radio

import (
	"context"
	"dubmer-bono/app/service"
	"dubmer-bono/app/types"

	"go.uber.org/zap"
)

type Service struct {
	repo              *service.Repository
	logger            *zap.SugaredLogger
	getDriverPressure func() int
}

func NewService(ctx context.Context, logger *zap.SugaredLogger, root string, repo *service.Repository, driver_pressure func() int) (types.Radio, error) {
	serv := &Service{repo: repo, logger: logger, getDriverPressure: driver_pressure}
	return serv, nil
}
