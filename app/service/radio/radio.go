package radio

import (
	"context"
	"dubmer-bono/app/service"
	"dubmer-bono/app/service/internal/audio/stt"
	"dubmer-bono/app/service/internal/audio/tts"
	"dubmer-bono/app/types"
	"dubmer-bono/app/types/entity"
	"dubmer-bono/app/utility"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	repo              *service.Repository
	logger            *zap.SugaredLogger
	getDriverPressure func() int
	prio_sorted_vc    map[int]*utility.ExpiryQueue[entity.RadioMessage]
	tts               *tts.TTS
	stt               *stt.STT
	msg_chan          chan entity.RadioMessage
	hotkey_chan       <-chan entity.HotKeyEvent
	muted             bool
	maplock           sync.RWMutex
	systemPrompt      string
	llmModel          string
}

const (
	defaultSystemPrompt = "You are the in-cab radio assistant for a driver. " +
		"Respond concisely and clearly, since your replies are spoken aloud over the radio. " +
		"Avoid long-winded answers, formatting, or anything that doesn't make sense read out loud." +
		"No Data source is connected to make shit up for now."

	defaultLLMModel = "llama-3.3-70b-versatile"
)

func NewService(ctx context.Context, logger *zap.SugaredLogger, root string, repo *service.Repository, driver_pressure func() int, msgchan chan entity.RadioMessage, hkeychan <-chan entity.HotKeyEvent) (types.Radio, error) {
	tts, err := tts.NewTTS(ctx)
	if err != nil {
		return Service{}, err
	}
	stt, err := stt.NewSTTHandler(ctx)
	if err != nil {
		return Service{}, err
	}
	serv := &Service{
		repo:              repo,
		logger:            logger,
		getDriverPressure: driver_pressure,
		tts:               tts,
		stt:               stt,
		msg_chan:          msgchan,
		prio_sorted_vc:    make(map[int]*utility.ExpiryQueue[entity.RadioMessage]),
		hotkey_chan:       hkeychan,
		systemPrompt:      defaultSystemPrompt,
		llmModel:          defaultLLMModel,
	}
	go serv.MessageChanListner()    // Listens to shared channel and updates queue
	go serv.radioTheDriver(ctx)     // Periodically checks queue and radios message when allowed
	go serv.HandleHotKeyEvents(ctx) //Act acc to the hot key events
	serv.stt.InitSTT(ctx)
	return serv, nil
}

// sentenceEndings are the characters that trigger a flush of the buffered
// text to TTS. Adjust as needed (e.g. add newline if your source emits one
// per sentence).
const sentenceEndings = ".!?"

func (s *Service) radioTheDriver(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			pressure := s.getDriverPressure()
			reader, ok := s.GetMessageByMinPriority()
			if !ok {
				continue
			}
			s.speakStream(ctx, reader, pressure)
		}
	}
}
