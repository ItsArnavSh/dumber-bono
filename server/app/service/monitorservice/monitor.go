package monitor

import (
	"context"
	"dubmer-bono/app/service"
	"dubmer-bono/app/types"
	"dubmer-bono/app/types/entity"

	"go.uber.org/zap"
)

/*
 * The monitor service is concerned with primarily listening for events
 *  And based on that, generating text, and pusing to the Radio pipeline
 */

type Service struct {
	repo            *service.Repository
	logger          *zap.SugaredLogger
	driver_pressure int
	msg_chan        chan<- entity.RadioMessage
}

var _ types.Monitor = &Service{}

func NewService(ctx context.Context, logger *zap.SugaredLogger, root string, repo *service.Repository, msg_chan chan<- entity.RadioMessage) (types.Monitor, error) {
	serv := &Service{repo: repo, logger: logger, msg_chan: msg_chan}
	go serv.monitorPressure(ctx)
	return serv, nil
}

func (s *Service) GetPressure() int {
	return s.driver_pressure
}
