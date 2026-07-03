package stt

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

type STT struct {
	URL string
}

func NewSTT(url string) *STT {
	return &STT{URL: url}
}

func (s *STT) StartLiveSTT(ctx context.Context, audioDispatcher <-chan []byte, textReceiver chan<- []byte) error {
	apiURL := s.URL + "/inference"

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case chunk, ok := <-audioDispatcher:
			if !ok {
				return nil
			}

			// Fire-and-forget the request in a goroutine
			// This keeps the loop free to receive more audio immediately
			go func(data []byte) {
				req, err := http.NewRequest("POST", apiURL, bytes.NewReader(data))
				if err != nil {
					return
				}
				req.Header.Set("Content-Type", "application/octet-stream")

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					return
				}
				defer resp.Body.Close()

				body, _ := io.ReadAll(resp.Body)
				if len(body) > 0 {
					textReceiver <- body
				}
			}(chunk)
		}
	}
}
