package audio

import (
	"context"
	"dubmer-bono/app/service/internal/audio"
	"fmt"
	"log"

	"github.com/gordonklaus/portaudio"
)

func (a *STT) recordAudio(ctx context.Context) {
	samples := make([]int16, audio.ChunkSize)
	stream := a.openAudioStream(samples)
	if stream == nil {
		return
	}
	defer stream.Close()

	if err := stream.Start(); err != nil {
		log.Fatal(err)
	}
	defer stream.Stop()

	fmt.Println("Recording...")
	a.processAudioLoop(stream, samples)
}

func (a *STT) openAudioStream(samples []int16) *portaudio.Stream {
	stream, err := portaudio.OpenDefaultStream(
		audio.Channels,
		0,
		float64(audio.SampleRate),
		len(samples),
		samples,
	)
	if err != nil {
		log.Fatal(err)
	}
	return stream
}

func (a *STT) processAudioLoop(stream *portaudio.Stream, samples []int16) {
	var (
		state     audioState
		pcmBuffer []int16
	)

	for {
		stream.Read()
		frame := a.copyFrame(samples)

		if a.isSpeechDetected(frame) {
			state.handleSpeech(&pcmBuffer, frame)
		} else if state.recording {
			state.handleSilence(&pcmBuffer, frame)
			if state.isSilenceThresholdReached() {
				a.finalizeRecording(&pcmBuffer, &state)
			}
		}
	}
}

type audioState struct {
	recording     bool
	silenceFrames int
}

func (a *STT) copyFrame(samples []int16) []int16 {
	frame := make([]int16, len(samples))
	copy(frame, samples)
	return frame
}

func (a *STT) isSpeechDetected(frame []int16) bool {
	speech, _ := a.vad.Process(16000, pcm16ToBytes(frame))
	return speech && (RMS(frame) > float64(Threshold))
}

func (s *audioState) handleSpeech(buffer *[]int16, frame []int16) {
	s.recording = true
	s.silenceFrames = 0
	*buffer = append(*buffer, frame...)
}

func (s *audioState) handleSilence(buffer *[]int16, frame []int16) {
	*buffer = append(*buffer, frame...)
	s.silenceFrames++
}

func (s *audioState) isSilenceThresholdReached() bool {
	return s.silenceFrames >= 25 // 25×20ms = 500ms
}

func (a *STT) finalizeRecording(buffer *[]int16, state *audioState) {
	const keepTrailing = 3 // 60 ms
	trimFrames := state.silenceFrames - keepTrailing

	if trimFrames > 0 {
		trimSamples := trimFrames * audio.ChunkSize
		if len(*buffer) >= trimSamples {
			*buffer = (*buffer)[:len(*buffer)-trimSamples]
		}
	}

	if wav, err := PCM16ToWAV(*buffer, 16000, 1); err == nil {
		{
			a.mu.Lock()
			if a.closed == false {
				a.incoming <- wav
			}
			a.mu.Unlock()
		}
	}

	*buffer = (*buffer)[:0]
	state.recording = false
	state.silenceFrames = 0
}
func (a *STT) Destructor() {
	portaudio.Terminate()
}
func pcm16ToBytes(samples []int16) []byte {
	b := make([]byte, len(samples)*2)

	for i, s := range samples {
		b[2*i] = byte(s)
		b[2*i+1] = byte(s >> 8)
	}

	return b
}
