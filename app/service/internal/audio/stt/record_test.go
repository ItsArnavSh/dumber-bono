package stt

import (
	"testing"
)

func TestAudioStateInitial(t *testing.T) {
	var s audioState
	if s.recording {
		t.Error("initial recording should be false")
	}
	if s.silenceFrames != 0 {
		t.Error("initial silenceFrames should be 0")
	}
}

func TestAudioStateHandleSpeech(t *testing.T) {
	var s audioState
	buffer := []int16{}
	frame := []int16{1, 2, 3}

	s.handleSpeech(&buffer, frame)
	if !s.recording {
		t.Error("recording should be true after speech")
	}
	if s.silenceFrames != 0 {
		t.Error("silenceFrames should reset to 0 after speech")
	}
	if len(buffer) != 3 {
		t.Errorf("buffer len = %d, want 3", len(buffer))
	}
}

func TestAudioStateHandleSilence(t *testing.T) {
	var s audioState
	buffer := []int16{1}
	frame := []int16{9, 9}

	s.handleSilence(&buffer, frame)
	if s.silenceFrames != 1 {
		t.Errorf("silenceFrames = %d, want 1", s.silenceFrames)
	}
	if len(buffer) != 3 {
		t.Errorf("buffer len = %d, want 3", len(buffer))
	}
}

func TestAudioStateSilenceThreshold(t *testing.T) {
	var s audioState
	if s.isSilenceThresholdReached() {
		t.Error("should not be at threshold initially")
	}

	for range 25 {
		if s.isSilenceThresholdReached() {
			break
		}
		s.silenceFrames++
	}
	if !s.isSilenceThresholdReached() {
		t.Error("should reach threshold after 25 frames")
	}
}

func TestCopyFrameDeepCopy(t *testing.T) {
	a := &STT{}
	samples := []int16{1, 2, 3}
	frame := a.copyFrame(samples)
	samples[0] = 999

	if frame[0] != 1 {
		t.Errorf("frame[0] = %d, want 1 (must be a deep copy)", frame[0])
	}
}

func TestFinalizeRecordingTrimsSilence(t *testing.T) {
	a := &STT{}
	a.incoming = make(chan []byte, 1)
	a.closed = false

	buffer := make([]int16, 30000) // ~1.8s of audio at 16kHz
	state := &audioState{recording: true, silenceFrames: 10}
	a.finalizeRecording(&buffer, state)

	if len(buffer) != 0 {
		t.Errorf("buffer len = %d, want 0 after finalize", len(buffer))
	}
	if state.recording {
		t.Error("state.recording should be false after finalize")
	}
}
