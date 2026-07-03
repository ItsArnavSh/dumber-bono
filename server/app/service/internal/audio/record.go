package audio

import (
	"context"
	"fmt"
	"log"

	"github.com/gordonklaus/portaudio"
)

func (a *Audio) RecordAudio(ctx context.Context) {
	samples := make([]int16, ChunkSize)

	stream, err := portaudio.OpenDefaultStream(
		Channels, // input channels
		0,        // output channels
		float64(SampleRate),
		len(samples),
		samples,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()
	if err := stream.Start(); err != nil {
		log.Fatal(err)
	}
	defer stream.Stop()

	var (
		recording     bool
		silenceFrames int
		pcmBuffer     []int16
	)
	fmt.Println("Recording...")
	for {
		stream.Read()

		frame := make([]int16, len(samples))
		copy(frame, samples)

		speech, _ := a.vad.Process(16000, PCM16ToBytes(frame))
		speech = speech && (RMS(frame) > float64(Threshold))
		switch {
		case speech:
			recording = true
			silenceFrames = 0
			pcmBuffer = append(pcmBuffer, frame...)

		case recording:
			// keep short pauses
			pcmBuffer = append(pcmBuffer, frame...)
			silenceFrames++

			if silenceFrames >= 25 { // 25×20ms = 500ms
				keepTrailing := 3 // 60 ms

				trimFrames := silenceFrames - keepTrailing
				if trimFrames > 0 {
					trimSamples := trimFrames * ChunkSize
					pcmBuffer = pcmBuffer[:len(pcmBuffer)-trimSamples]
				}
				wav, err := PCM16ToWAV(pcmBuffer, 16000, 1)
				if err == nil {
					a.incoming <- wav
				}

				pcmBuffer = pcmBuffer[:0]
				recording = false
				silenceFrames = 0
			}
		}
	}
}
func processChunk(samples []int16) {
	fmt.Println("First sample:", samples[0])
}
func (a *Audio) Destructor() {
	portaudio.Terminate()
}
func PCM16ToBytes(samples []int16) []byte {
	b := make([]byte, len(samples)*2)

	for i, s := range samples {
		b[2*i] = byte(s)
		b[2*i+1] = byte(s >> 8)
	}

	return b
}
