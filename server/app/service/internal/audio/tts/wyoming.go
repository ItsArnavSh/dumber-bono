package tts

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
)

type wyomingHeader struct {
	Type          string          `json:"type"`
	Version       string          `json:"version,omitempty"`
	Data          json.RawMessage `json:"data,omitempty"`
	DataLength    *int            `json:"data_length,omitempty"`
	PayloadLength *int            `json:"payload_length,omitempty"`
}

type synthesizeData struct {
	Text  string     `json:"text"`
	Voice *voiceData `json:"voice,omitempty"`
}

type voiceData struct {
	Name string `json:"name,omitempty"`
}

type audioStartData struct {
	Rate     int `json:"rate"`
	Width    int `json:"width"`
	Channels int `json:"channels"`
}

type wyomingConn struct {
	conn net.Conn
	br   *bufio.Reader
}

func newWyomingConn(conn net.Conn) *wyomingConn {
	return &wyomingConn{
		conn: conn,
		br:   bufio.NewReader(conn),
	}
}

func (w *wyomingConn) writeMessage(msgType string, data any, payload []byte) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal data: %w", err)
	}

	hdr := wyomingHeader{Type: msgType, Data: dataBytes}
	if payload != nil {
		n := len(payload)
		hdr.PayloadLength = &n
	}

	hdrBytes, err := json.Marshal(hdr)
	if err != nil {
		return fmt.Errorf("marshal header: %w", err)
	}

	if _, err := w.conn.Write(append(hdrBytes, '\n')); err != nil {
		return fmt.Errorf("write header: %w", err)
	}
	if payload != nil {
		if _, err := w.conn.Write(payload); err != nil {
			return fmt.Errorf("write payload: %w", err)
		}
	}
	return nil
}

// readMessage reads one newline-terminated header line. If data_length is
// present, it then reads that many bytes as a separate JSON metadata block
// (this is how wyoming-piper 1.8.0 actually sends "data" — NOT inline in
// the header). If payload_length is present, it reads that many raw bytes
// as the binary payload.
func (w *wyomingConn) readMessage() (wyomingHeader, []byte, error) {
	line, err := w.br.ReadBytes('\n')
	if err != nil {
		return wyomingHeader{}, nil, fmt.Errorf("read header line: %w", err)
	}

	var hdr wyomingHeader
	if err := json.Unmarshal(line, &hdr); err != nil {
		return wyomingHeader{}, nil, fmt.Errorf("unmarshal header: %w (line=%q)", err, line)
	}

	if hdr.DataLength != nil && *hdr.DataLength > 0 {
		dataBytes := make([]byte, *hdr.DataLength)
		if _, err := io.ReadFull(w.br, dataBytes); err != nil {
			return hdr, nil, fmt.Errorf("read data block: %w", err)
		}
		hdr.Data = dataBytes
	}

	if hdr.PayloadLength == nil || *hdr.PayloadLength <= 0 {
		return hdr, nil, nil
	}

	payload := make([]byte, *hdr.PayloadLength)
	if _, err := io.ReadFull(w.br, payload); err != nil {
		return hdr, nil, fmt.Errorf("read payload: %w", err)
	}

	return hdr, payload, nil
}

func pcmBytesToInt16(b []byte) []int16 {
	n := len(b) / 2
	out := make([]int16, n)
	for i := range n {
		out[i] = int16(binary.LittleEndian.Uint16(b[i*2 : i*2+2]))
	}
	return out
}
func (t *TTS) synthesizeOne(ctx context.Context, wc *wyomingConn, sentence string, w io.Writer) error {
	if err := wc.writeMessage("synthesize", synthesizeData{Text: sentence}, nil); err != nil {
		return fmt.Errorf("send synthesize: %w", err)
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		hdr, payload, err := wc.readMessage()
		if err != nil {
			return err
		}
		switch hdr.Type {
		case "audio-start":
			var start audioStartData
			_ = json.Unmarshal(hdr.Data, &start)
			fmt.Printf("[TTS] audio-start rate=%d channels=%d width=%d\n", start.Rate, start.Channels, start.Width)
		case "audio-chunk":
			if _, err := w.Write(payload); err != nil {
				return err
			}
		case "audio-stop":
			return nil
		case "error":
			return fmt.Errorf("server error: %s", string(hdr.Data))
		}
	}
}
