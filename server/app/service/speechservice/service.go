package speechservice

import (
	"context"
	"dubmer-bono/app/service/internal/audio"
	"dubmer-bono/app/service/internal/stt"
	"dubmer-bono/app/service/internal/tts"
	"fmt"
)

type Service struct {
	stt stt.STT
	tts tts.TTS
}

func StartListner(ctx context.Context) {
	receiver := make(chan []byte, 100) // Buffered to prevent dropping audio
	textchan := make(chan []byte, 100)

	// Start Audio
	go func() {
		if err := audio.Listen(ctx, receiver); err != nil {
			fmt.Printf("Audio error: %v\n", err)
		}
	}()

	// Start STT
	sttSvc := stt.NewSTT("ws://localhost:8081")
	go func() {
		if err := sttSvc.StartLiveSTT(ctx, receiver, textchan); err != nil {
			fmt.Printf("STT Connection error: %v\n", err)
		}
	}()

	fmt.Println("Started Listening...")

	// Consume transcripts
	for text := range textchan {
		// Convert the byte slice to a readable string
		fmt.Printf("Transcript: %s\n", string(text))
	}
}
