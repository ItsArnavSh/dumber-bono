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
}

func NewService(ctx context.Context, logger *zap.SugaredLogger, root string, repo *service.Repository, driver_pressure func() int) (types.Radio, error) {
	tts, err := tts.NewTTS()
	if err != nil {
		return Service{}, err
	}
	serv := &Service{
		repo:              repo,
		logger:            logger,
		getDriverPressure: driver_pressure,
		tts:               tts,
	}
	serv.TestSpeaker(ctx)
	return serv, nil
}

func (s *Service) TestSpeaker(ctx context.Context) {
	text := "Hi Hamilton, Radio Check. " +
		"Box box box, confirm undercut on Verstappen. " +
		"Piastri and Norris are side by side into Eau Rouge. " +
		"Leclerc, DRS enabled, gap to Russell is nine tenths. " +
		"Antonelli locking up into the chicane. " +
		"Alonso and Stroll running the medium compound. " +
		"Sainz and Albon on an alternate strategy at Williams. " +
		"Hulkenberg and Bortoleto for Audi, both on fresh softs. " +
		"Gasly and Colapinto battling for P8 at Alpine. " +
		"Ocon and Bearman holding position for Haas. " +
		"Perez and Bottas debuting for Cadillac this weekend. " +
		"Hadjar and Lawson through the final sector. " +
		"Lindblad on his out-lap for Racing Bulls. " +
		"Push now, push now, delta minus point-three. " +
		"Safety car deployed, virtual safety car ending. " +
		"Parc ferme conditions apply after qualifying."
	r, _ := s.tts.StringToPCM(ctx, text, 2)
	_ = speaker.PlayPCM(ctx, r)
}
