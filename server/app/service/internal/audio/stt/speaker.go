package audio

import (
	"context"
	"log"

	"github.com/gordonklaus/portaudio"
)

func (a *Audio) PlayPCM(ctx context.Context) error {
	const writeBlockSize = ChunkSize // write smaller, more frequent blocks
	samples := make([]int16, writeBlockSize)

	stream, err := portaudio.OpenDefaultStream(
		0,
		1,
		float64(22050),
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

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case chunk, ok := <-a.outgoing:
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

			buf = append(buf, chunk...)

			if err := writeLoop(); err != nil {
				log.Printf("stream.Write failed: %v", err)
				// Don't die on a transient underflow — log and keep going,
				// otherwise one glitch kills the whole utterance.
				continue
			}
		}
	}
}
