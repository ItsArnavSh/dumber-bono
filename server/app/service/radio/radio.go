package radio

import (
	"context"
	"dubmer-bono/app/service"
	"dubmer-bono/app/service/internal/audio/speaker"
	"dubmer-bono/app/service/internal/audio/tts"
	"dubmer-bono/app/types"
	"dubmer-bono/app/types/entity"
	"dubmer-bono/app/utility"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	repo              *service.Repository
	logger            *zap.SugaredLogger
	getDriverPressure func() int
	prio_sorted_vc    map[int]*utility.ExpiryQueue[entity.RadioMessage]
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
		prio_sorted_vc:    make(map[int]*utility.ExpiryQueue[entity.RadioMessage]),
	}
	go serv.MessageChanListner() // Listens to shared channel and updates queue
	go serv.radioTheDriver(ctx)  // Periodically checks queue and radios message when allowed
	return serv, nil
}

func (s *Service) radioTheDriver(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			pressure := s.getDriverPressure()

			message, ok := s.GetMessageByMinPriority()
			if !ok {
				continue
			}

			pcm, err := s.tts.StringToPCM(ctx, message)
			if err != nil {
				s.logger.Errorf("error converting string to PCM: %v", err)
				continue
			}

			s.logger.Infof("Engineer: %s \n Priority: %d", message, pressure)

			if err := speaker.PlayPCM(ctx, pcm); err != nil {
				s.logger.Errorf("error playing PCM: %v", err)
				continue
			}
			s.logger.Debugf("finished playing message at priority %d", pressure)
		}
	}
}
