package audio

import (
	"context"
	"fmt"

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
}

func NewAudioHandler(ctx context.Context, incoming, outgoing chan []byte) (Audio, error) {
	vad, err := webrtcvad.New()
	if err != nil {
		return Audio{}, err
	}
	vad.SetMode(3)
	return Audio{
		incoming: incoming,
		outgoing: outgoing,
		vad:      vad,
	}, nil
}

func (a *Audio) SetupPortAudio(ctx context.Context) error {
	err := portaudio.Initialize()
	return err
}

func (a *Audio) ProcessWAVChunks(ctx context.Context) {
	for _ = range a.incoming {
		fmt.Println("Received the Speech Chunk")
	}
}
