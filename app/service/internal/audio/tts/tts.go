package tts

import (
	"context"
	"io"
	"net"
)

type TTS struct {
	wc      *wyomingConn
	rawConn net.Conn
}

func NewTTS(ctx context.Context) (*TTS, error) {
	dialer := net.Dialer{}
	rawConn, err := dialer.DialContext(ctx, "tcp", "127.0.0.1:10200")
	if err != nil {
		return nil, err
	}
	wc := newWyomingConn(rawConn)
	return &TTS{wc: wc, rawConn: rawConn}, nil
}

func (t *TTS) Close() error {
	return t.rawConn.Close()
}

func (t *TTS) StringToPCM(ctx context.Context, sentence string) (io.Reader, error) {
	pr, pw := io.Pipe()
	go func() {
		if err := t.synthesizeOne(ctx, t.wc, sentence, pw); err != nil {
			pw.CloseWithError(err)
			return
		}
		_ = pw.Close()
	}()
	return pr, nil
}
