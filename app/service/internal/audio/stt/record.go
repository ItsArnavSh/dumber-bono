package stt

import (
	"context"
	"dubmer-bono/app/service/internal/audio"
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
		if err := stream.Read(); err != nil {
			continue
		}

		a.mu.Lock()
		active := a.isRecording
		a.mu.Unlock()

		// If recording stopped, force-finalize any leftover speech immediately
		if !active {
			if state.recording && len(pcmBuffer) > 0 {
				log.Println("[STT FORCE FLUSH] EndMessageRec called mid-speech. Force-finalizing buffer...")
				a.finalizeRecording(&pcmBuffer, &state)
			}
			continue
		}

		frame := a.copyFrame(samples)
		isSpeech := a.isSpeechDetected(frame)

		if isSpeech {
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
	speech, err := a.vad.Process(16000, pcm16ToBytes(frame))
	if err != nil {
		log.Printf("[VAD ERROR] Processing failed: %v", err)
	}

	rmsVal := RMS(frame)
	thresholdVal := float64(Threshold)
	detected := speech && (rmsVal > thresholdVal)

	// // Logging raw energy level and VAD state to trace noise spikes
	// if speech || rmsVal > (thresholdVal*0.5) {
	// 	log.Printf("[VAD METRICS] VAD Speech: %v | RMS: %.2f (Threshold: %.2f) -> Active: %v", speech, rmsVal, thresholdVal, detected)
	// }

	return detected
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

	totalSamples := len(*buffer)

	// 1. Trim trailing silence frames from the buffer
	if trimFrames > 0 {
		trimSamples := trimFrames * audio.ChunkSize
		if totalSamples >= trimSamples {
			*buffer = (*buffer)[:totalSamples-trimSamples]
		}
	}

	trimmedSamples := len(*buffer)
	log.Printf("[STT FINALIZE] Captured %d raw samples (Trimmed to %d samples, ~%.2fs)",
		totalSamples, trimmedSamples, float64(trimmedSamples)/16000.0)

	// 2. Minimum Duration Guard
	const minSamples = 12800 // 0.8s * 16,000 samples/sec
	if trimmedSamples < minSamples {
		log.Printf("[STT DROPPED] Audio chunk ignored: duration too short (%.2fs < 0.80s threshold)",
			float64(trimmedSamples)/16000.0)
	} else {
		// 3. Deep-copy buffer memory before passing to PCM16ToWAV
		audioCopy := make([]int16, trimmedSamples)
		copy(audioCopy, *buffer)

		wav, err := PCM16ToWAV(audioCopy, 16000, 1)
		if err != nil {
			log.Printf("[STT ERROR] Failed to encode PCM to WAV: %v", err)
		} else {
			a.mu.Lock()
			if a.closed {
				log.Println("[STT DROPPED] Audio ready but channel 'a.incoming' is closed")
			} else {
				log.Printf("[STT QUEUED] Pushing %d bytes WAV to channel 'a.incoming'...", len(wav))

				// Non-blocking write check to diagnose stuck consumers
				select {
				case a.incoming <- wav:
					log.Println("[STT SENT] WAV payload successfully delivered to transcriber")
				default:
					log.Println("[STT WARN] 'a.incoming' channel buffer full or no active receiver loop!")
					// Fallback to blocking send if necessary
					a.incoming <- wav
				}
			}
			a.mu.Unlock()
		}
	}

	// 4. Force allocation of a new backing array on next recording cycle
	*buffer = nil
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
