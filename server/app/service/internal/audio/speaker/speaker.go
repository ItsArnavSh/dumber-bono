package speaker

import (
	"context"
	"dubmer-bono/app/service/internal/audio"
	"encoding/binary"
	"io"
	"log"

	"github.com/gordonklaus/portaudio"
)

func PlayPCM(ctx context.Context, r io.Reader) error {
	const writeBlockSize = audio.ChunkSize // write smaller, more frequent blocks
	samples := make([]int16, writeBlockSize)
	stream, err := portaudio.OpenDefaultStream(
		0,
		1,
		float64(22050)*1.1,
		len(samples),
		samples,
	)
	if err != nil {
		log.Printf("failed to open output stream: %v", err)
		return err
	}
	defer func() {
		if err := stream.Close(); err != nil {
			log.Printf("failed to close output stream: %v", err)
		}
	}()
	if err := stream.Start(); err != nil {
		log.Printf("failed to start output stream: %v", err)
		return err
	}
	defer func() {
		if err := stream.Stop(); err != nil {
			log.Printf("failed to stop output stream: %v", err)
		}
	}()
	log.Println("Playback started")

	var buf []int16
	writeLoop := func() error {
		for len(buf) >= writeBlockSize {
			copy(samples, buf[:writeBlockSize])
			buf = buf[writeBlockSize:]
			if err := stream.Write(); err != nil {
				return err
			}
		}
		return nil
	}

	// Read raw PCM bytes off r in a goroutine, decode to int16, push into a channel
	// so we can still select on ctx.Done() while reading.
	type readResult struct {
		samples []int16
		err     error
	}
	results := make(chan readResult)
	go func() {
		defer close(results)
		raw := make([]byte, writeBlockSize*2) // 2 bytes per int16 sample
		for {
			n, err := io.ReadFull(r, raw)
			if n > 0 {
				// n may be < len(raw) on the final partial read (io.ErrUnexpectedEOF)
				full := n - (n % 2) // drop a dangling odd byte, if any
				out := make([]int16, full/2)
				for i := range out {
					out[i] = int16(binary.LittleEndian.Uint16(raw[i*2:]))
				}
				select {
				case results <- readResult{samples: out}:
				case <-ctx.Done():
					return
				}
			}
			if err != nil {
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					return
				}
				select {
				case results <- readResult{err: err}:
				case <-ctx.Done():
				}
				return
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case res, ok := <-results:
			if !ok {
				// Final flush, padded.
				if len(buf) > 0 {
					n := copy(samples, buf)
					for i := n; i < len(samples); i++ {
						samples[i] = 0
					}
					buf = nil
					_ = stream.Write()
				}
				return nil
			}
			if res.err != nil {
				log.Printf("PCM read failed: %v", res.err)
				return res.err
			}
			buf = append(buf, res.samples...)
			if err := writeLoop(); err != nil {
				log.Printf("stream.Write failed: %v", err)
				// Don't die on a transient underflow — log and keep going,
				// otherwise one glitch kills the whole utterance.
				continue
			}
		}
	}
}
