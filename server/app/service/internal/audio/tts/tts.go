package tts

import (
	"context"
	"io"
	"net"
	"sync"
)

type TTS struct {
	wc      *wyomingConn
	rawConn net.Conn
	mu      sync.Mutex // guards synthesizeOne against concurrent calls
}

func NewTTS() (*TTS, error) {
	rawConn, err := net.Dial("tcp", "127.0.0.1:10200")
	if err != nil {
		return nil, err
	}
	wc := newWyomingConn(rawConn)
	return &TTS{wc: wc, rawConn: rawConn}, nil
}

func (t *TTS) Close() error {
	return t.rawConn.Close()
}

func (t *TTS) StringToPCM(ctx context.Context, sentence string, speed float64) (io.Reader, error) {
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		if err := t.synthesizeOne(ctx, t.wc, sentence, pw); err != nil {
			pw.CloseWithError(err)
		}
	}()
	return pr, nil
}
