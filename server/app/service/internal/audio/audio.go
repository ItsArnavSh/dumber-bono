package audio

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/gordonklaus/portaudio"
)

// Listen starts an audio stream and pipes chunks into the receiver channel.
// It will block until the context is cancelled.
func Listen(ctx context.Context, receiver chan<- []byte) error {
	err := portaudio.Initialize()
	if err != nil {
		return fmt.Errorf("failed to initialize portaudio: %w", err)
	}
	// Ensure termination happens when the function exits
	defer portaudio.Terminate()

	const bufferSize = 512
	inputBuffer := make([]int16, bufferSize)

	// 1 input channel, 0 output, 16kHz sample rate
	stream, err := portaudio.OpenDefaultStream(1, 0, 16000, bufferSize, inputBuffer)
	if err != nil {
		return fmt.Errorf("failed to open stream: %w", err)
	}
	defer stream.Close()

	if err := stream.Start(); err != nil {
		return fmt.Errorf("failed to start stream: %w", err)
	}
	defer stream.Stop()

	fmt.Println("Audio listener started...")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Audio listener stopping...")
			return ctx.Err()
		default:
			// Read() blocks until buffer is full
			if err := stream.Read(); err != nil {
				fmt.Printf("Read error: %v\n", err)
				continue
			}

			// Send a copy of the buffer to avoid data races
			// since we are reusing inputBuffer in the next iteration
			chunk := make([]byte, len(inputBuffer)*2)
			for i, v := range inputBuffer {
				binary.LittleEndian.PutUint16(chunk[i*2:], uint16(v))
			}
			receiver <- chunk
		}
	}
}

//Use Vosk Instead
