package agent

import (
	"context"
)

type Service struct {
}

func StartListner(ctx context.Context) {
	// incoming := make(chan []byte, 100) // Buffered to prevent dropping audio
	// outgoing := make(chan []int16, 100)
	// textchan := make(chan string, 100)

	// ahandler, _ := audio.NewAudioHandler(ctx, incoming, outgoing, textchan)
	// ahandler.SetupPortAudio(ctx)
	// go ahandler.TranscribeWAV(ctx)
	// go ahandler.RecordAudio(ctx)
	// go ahandler.TTS(ctx)
	// go ahandler.PlayPCM(ctx)
	// select {}
}
