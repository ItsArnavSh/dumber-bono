package stt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// STT talks to a Vosk server (e.g. alphacep/kaldi-en) over its
// WebSocket protocol: https://github.com/alphacep/vosk-server
type STT struct {
	// URL should be a ws:// address, e.g. "ws://localhost:2700"
	URL        string
	SampleRate int // e.g. 16000, must match the audio you send
}

func NewSTT(url string) *STT {
	return &STT{URL: url, SampleRate: 16000}
}

// voskResult models both partial and final Vosk responses.
// Only one of Text/Partial will be non-empty per message.
type voskResult struct {
	Text    string `json:"text,omitempty"`
	Partial string `json:"partial,omitempty"`
}

// StartLiveSTT opens a single websocket connection to Vosk, streams every
// audio chunk from audioDispatcher into it, and forwards recognized text
// (final results only) to textReceiver. It blocks until ctx is cancelled
// or audioDispatcher is closed.
func (s *STT) StartLiveSTT(ctx context.Context, audioDispatcher <-chan []byte, textReceiver chan<- []byte) error {
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, s.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to vosk server: %w", err)
	}
	defer conn.Close()

	// Tell Vosk the sample rate of the audio we're about to stream.
	cfg := fmt.Sprintf(`{"config": {"sample_rate": %d}}`, s.SampleRate)
	if err := conn.WriteMessage(websocket.TextMessage, []byte(cfg)); err != nil {
		return fmt.Errorf("failed to send vosk config: %w", err)
	}

	// Reader goroutine: pulls recognition results off the socket as they
	// arrive, independent of when we happen to be writing audio.
	readerDone := make(chan struct{})
	go func() {
		defer close(readerDone)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return // socket closed, e.g. after we send eof or ctx is done
			}
			var res voskResult
			if err := json.Unmarshal(message, &res); err != nil {
				log.Printf("stt: bad json from vosk: %v", err)
				continue
			}
			// Only forward finalized text. Drop res.Partial silently;
			// wire it up here too if you want live partial captions.
			if res.Text == "" {
				continue
			}
			select {
			case textReceiver <- []byte(res.Text):
			case <-ctx.Done():
				return
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			s.closeStream(conn, readerDone)
			return ctx.Err()

		case chunk, ok := <-audioDispatcher:
			if !ok {
				s.closeStream(conn, readerDone)
				return nil
			}
			if err := conn.WriteMessage(websocket.BinaryMessage, chunk); err != nil {
				return fmt.Errorf("failed to write audio chunk to vosk: %w", err)
			}
		}
	}
}

// closeStream tells Vosk we're done sending audio (which flushes a final
// result for anything still buffered) and waits for the reader goroutine
// to notice the socket close.
func (s *STT) closeStream(conn *websocket.Conn, readerDone <-chan struct{}) {
	_ = conn.WriteMessage(websocket.TextMessage, []byte(`{"eof" : 1}`))
	<-readerDone
}
