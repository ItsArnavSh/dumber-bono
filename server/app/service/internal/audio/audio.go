package audio

import (
	"context"

	"github.com/gordonklaus/portaudio"
	"github.com/maxhawkins/go-webrtcvad"
)

const (
	SampleRate = 16000
	Channels   = 1
	ChunkSize  = 320 // 20ms
)

type Audio struct {
	incoming chan []byte
	outgoing chan []byte
	vad      *webrtcvad.VAD
	whisper  *Whisper
}

func NewAudioHandler(ctx context.Context, incoming, outgoing chan []byte) (Audio, error) {
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
		incoming: incoming,
		outgoing: outgoing,
		vad:      vad,
		whisper:  whisper,
	}, nil
}

func (a *Audio) SetupPortAudio(ctx context.Context) error {
	err := portaudio.Initialize()
	return err
}
