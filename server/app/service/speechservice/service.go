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
	outgoing := make(chan []int16, 100)
	textchan := make(chan string, 100)

	ahandler, _ := audio.NewAudioHandler(ctx, incoming, outgoing, textchan)
	ahandler.SetupPortAudio(ctx)
	go ahandler.TranscribeWAV(ctx)
	go ahandler.RecordAudio(ctx)
	go ahandler.TTS(ctx)
	go ahandler.PlayPCM(ctx)
	select {}
}
