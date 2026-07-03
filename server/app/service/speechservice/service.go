package speechservice

import (
	"context"
	"dubmer-bono/app/service/internal/audio"
	"dubmer-bono/app/service/internal/stt"
	"dubmer-bono/app/service/internal/tts"
)

type Service struct {
	stt stt.STT
	tts tts.TTS
}

func StartListner(ctx context.Context) {
	incoming := make(chan []byte, 100) // Buffered to prevent dropping audio
	outgoing := make(chan []byte, 100)

	ahandler, _ := audio.NewAudioHandler(ctx, incoming, outgoing)
	ahandler.SetupPortAudio(ctx)
	go ahandler.ProcessWAVChunks(ctx)
	ahandler.RecordAudio(ctx)
}
