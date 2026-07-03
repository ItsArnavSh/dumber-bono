package tts

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"net"
)

type TTS struct {
	Addr string
}

func NewTTS(addr string) *TTS {
	return &TTS{Addr: addr}
}

func (t *TTS) StartStreamAudio(ctx context.Context, textChan chan []byte) (chan []byte, error) {
	conn, err := net.Dial("tcp", t.Addr)
	if err != nil {
		return nil, err
	}

	audioChan := make(chan []byte)

	go func() {
		defer conn.Close()
		defer close(audioChan)

		for {
			select {
			case <-ctx.Done():
				return
			case text, ok := <-textChan:
				if !ok {
					return
				}

				// 1. Send Wyoming 'synthesize' event
				msg := map[string]any{"text": string(text)}
				data, _ := json.Marshal(msg)

				// Header: 4 bytes length, then JSON
				header := make([]byte, 4)
				binary.LittleEndian.PutUint32(header, uint32(len(data)))
				conn.Write(header)
				conn.Write(data)

				// 2. Read back audio chunks
				// (Simplified logic: in a real implementation,
				// you'd loop until the 'audio-stop' event)
				buf := make([]byte, 1024)
				n, _ := conn.Read(buf)
				audioChan <- buf[:n]
			}
		}
	}()

	return audioChan, nil
}
