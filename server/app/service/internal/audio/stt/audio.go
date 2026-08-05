package stt

import (
	"context"
	"strings"
	"sync"

	"github.com/maxhawkins/go-webrtcvad"
)

type STT struct {
	incoming       chan []byte
	vad            *webrtcvad.VAD
	whisper        *Whisper
	StrChunk       *strings.Builder
	SessionMessage string
	mu             sync.Mutex
	closed         bool //Flag to avoid sending in a closed channel

	done chan error // Signal channel for completion
	wg   sync.WaitGroup
}

func NewSTTHandler(ctx context.Context) (*STT, error) {
	incoming := make(chan []byte)
	vad, err := webrtcvad.New()
	if err != nil {
		return &STT{}, err
	}
	vad.SetMode(3)
	whisper := newWhisper()
	return &STT{
		incoming: incoming,
		vad:      vad,
		whisper:  whisper,
	}, nil
}

func (a *STT) InitSTT(ctx context.Context) {
	go a.recordAudio(ctx)
}

func (a *STT) StartMessageRec(ctx context.Context) {
	// Initialize the builder
	a.StrChunk = &strings.Builder{}

	// Open Incoming
	{
		a.mu.Lock()
		a.incoming = make(chan []byte)
		a.done = make(chan error, 1) // Buffered to avoid goroutine leak
		a.closed = false
		a.mu.Unlock()
	}

	// Start transcription in parallel
	a.wg.Add(1)
	go a.transcribeWAV(ctx)
}

func (a *STT) EndMessageRec(ctx context.Context) (string, error) {
	// Close Incoming - signals transcribeWAV to stop iterating
	{
		a.mu.Lock()
		close(a.incoming)
		a.closed = true
		a.mu.Unlock()
	}
	// Wait for transcribeWAV goroutine to finish
	a.wg.Wait()

	// Get result from done channel
	err := <-a.done
	if err != nil {
		return "", err
	}

	return a.StrChunk.String(), nil
}
