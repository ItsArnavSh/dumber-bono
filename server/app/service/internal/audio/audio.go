package audio

import (
	"context"
	"dubmer-bono/app/service/internal/llm"
	"log"

	"github.com/gordonklaus/portaudio"
	"github.com/maxhawkins/go-webrtcvad"
)

const (
	SampleRate = 16000
	Channels   = 1
	ChunkSize  = 320 // 20ms
)

type Audio struct {
	incoming   chan []byte
	outgoing   chan []int16
	vad        *webrtcvad.VAD
	whisper    *Whisper
	wordchunks chan string
}

func NewAudioHandler(ctx context.Context, incoming chan []byte, outgoing chan []int16, wordchunks chan string) (Audio, error) {
	vad, err := webrtcvad.New()
	if err != nil {
		return Audio{}, err
	}
	vad.SetMode(3)
	whisper := newWhisper(
		"../bin/whisper-cli",
		"../models/ggml-base.en.bin",
	)
	return Audio{
		incoming:   incoming,
		outgoing:   outgoing,
		vad:        vad,
		whisper:    whisper,
		wordchunks: wordchunks,
	}, nil
}

func (a *Audio) SetupPortAudio(ctx context.Context) error {
	err := portaudio.Initialize()
	return err
}

func (a *Audio) InvokeLLM(ctx context.Context, prompt string) error {
	tokens := make(chan string, 32)
	go func() {
		systemPrompt := "You are an F1 race engineer. While testing, named Dumber Bono. Make shit up for now. Give conversational very very short answers Its a high panic situation you get a few seconds to speak. Like one liners"
		if err := llm.StreamLLM(ctx, systemPrompt, prompt, "openai/gpt-oss-20b", tokens); err != nil {
			log.Println("stream error:", err)
		}
	}()
	for tok := range tokens {
		a.wordchunks <- tok
	}
	return nil
}
