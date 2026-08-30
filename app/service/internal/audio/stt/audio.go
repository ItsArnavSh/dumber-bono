package stt

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/maxhawkins/go-webrtcvad"
)

type STT struct {
	incoming       chan []byte
	vad            *webrtcvad.VAD
	whisper        *Whisper
	StrChunk       *strings.Builder
	SessionMessage string
	mu             sync.Mutex
	closed         bool // Flag to avoid sending in a closed channel
	isRecording    bool
	done           chan error // Signal channel for completion
	wg             sync.WaitGroup
}

func NewSTTHandler(ctx context.Context) (*STT, error) {
	incoming := make(chan []byte, 10)
	vad, err := webrtcvad.New()
	if err != nil {
		return &STT{}, err
	}
	if err := vad.SetMode(3); err != nil {
		return &STT{}, err
	}
	whisper := newWhisper()

	log.Println("[STT INIT] STT Handler initialized")
	return &STT{
		incoming: incoming,
		vad:      vad,
		whisper:  whisper,
	}, nil
}

func (a *STT) InitSTT(ctx context.Context) {
	log.Println("[STT INIT] Starting background audio recording thread...")
	go a.recordAudio(ctx)
}

func (a *STT) StartMessageRec(ctx context.Context) {
	// Initialize the builder
	a.StrChunk = &strings.Builder{}

	// Open Incoming
	{
		a.mu.Lock()
		a.incoming = make(chan []byte, 10)
		a.done = make(chan error, 1) // Buffered to avoid goroutine leak
		a.closed = false
		a.isRecording = true // Enable audio processing
		log.Println("[STT STATE] Channel 'a.incoming' OPENED for new message session")
		a.mu.Unlock()
	}

	// Start transcription in parallel
	a.wg.Add(1)
	go a.transcribeWAV(ctx)
}

func (a *STT) EndMessageRec(ctx context.Context) (string, error) {
	a.mu.Lock()
	a.isRecording = false // Stop capturing immediately
	a.mu.Unlock()

	// Give processAudioLoop up to 100ms to finalize any buffered speech before closing channel
	time.Sleep(100 * time.Millisecond)

	a.mu.Lock()
	if !a.closed {
		close(a.incoming)
		a.closed = true
		log.Println("[STT STATE] Channel 'a.incoming' CLOSED - flushing remaining chunks to worker")
	}
	a.mu.Unlock()

	a.wg.Wait()

	err := <-a.done
	if err != nil {
		return "", err
	}

	finalText := a.StrChunk.String()
	log.Printf("[STT STATE COMPLETE] Message session ended. Total text length: %d chars", len(finalText))
	return finalText, nil
}
