package radio

import (
	"context"
	"dubmer-bono/app/service"
	"dubmer-bono/app/service/internal/audio/speaker"
	"dubmer-bono/app/service/internal/audio/tts"
	"dubmer-bono/app/types"
	"dubmer-bono/app/types/entity"
	"dubmer-bono/app/utility"

	"go.uber.org/zap"
)

type Service struct {
	repo              *service.Repository
	logger            *zap.SugaredLogger
	getDriverPressure func() int
	prio_sorted_vc    map[int]utility.ExpiryQueue[entity.RadioMessage]
	tts               *tts.TTS
	msg_chan          <-chan entity.RadioMessage
}

func NewService(ctx context.Context, logger *zap.SugaredLogger, root string, repo *service.Repository, driver_pressure func() int, msgchan <-chan entity.RadioMessage) (types.Radio, error) {
	tts, err := tts.NewTTS()
	if err != nil {
		return Service{}, err
	}
	serv := &Service{
		repo:              repo,
		logger:            logger,
		getDriverPressure: driver_pressure,
		tts:               tts,
		msg_chan:          msgchan,
	}
	return serv, nil
}

func (s *Service) radioDriver(ctx context.Context) {
	for {
		pressure := s.getDriverPressure()
		message, ok := s.GetMessageByMinPriority(pressure)
		if !ok {
			continue
		}
		r, err := s.tts.StringToPCM(ctx, message)
		if err != nil {
			s.logger.Errorf("Error converting string%w", err)
		}
		s.logger.Infof("Engineer: %s \n Priority: %d", message, pressure)
		_ = speaker.PlayPCM(ctx, r)
	}
}
