package stt

import (
	"math"
	"testing"
)

func TestRMSZeroSamples(t *testing.T) {
	if got := RMS([]int16{0, 0, 0, 0}); got != 0 {
		t.Fatalf("RMS = %v, want 0", got)
	}
}

func TestRMSConstantAmplitude(t *testing.T) {
	// RMS of [100,100,100] == 100
	if got := RMS([]int16{100, 100, 100}); got != 100 {
		t.Fatalf("RMS = %v, want 100", got)
	}
}

func TestRMSAlternatingSigns(t *testing.T) {
	// (-100)^2 and 100^2 both equal 10000, so RMS == 100
	got := RMS([]int16{-100, 100, -100, 100})
	if got != 100 {
		t.Fatalf("RMS = %v, want 100", got)
	}
}

func TestRMSKnownValue(t *testing.T) {
	samples := []int16{3, 4}
	// (9+16)/2 = 12.5, sqrt(12.5) ~= 3.5355
	got := RMS(samples)
	want := math.Sqrt(12.5)
	if math.Abs(got-want) > 1e-9 {
		t.Fatalf("RMS = %v, want %v", got, want)
	}
}

func TestRMSMonoToneIsNonNegative(t *testing.T) {
	if got := RMS([]int16{1, 2, 3}); got < 0 {
		t.Fatalf("RMS = %v, want >= 0", got)
	}
}

func TestPCM16ToWAVHeader(t *testing.T) {
	samples := []int16{100, -100, 0}
	wav, err := PCM16ToWAV(samples, 16000, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// RIFF/WAVE magic
	if string(wav[0:4]) != "RIFF" {
		t.Fatalf("missing RIFF header, got %q", wav[0:4])
	}
	if string(wav[8:12]) != "WAVE" {
		t.Fatalf("missing WAVE marker, got %q", wav[8:12])
	}
	if string(wav[12:16]) != "fmt " {
		t.Fatalf("missing fmt chunk, got %q", wav[12:16])
	}
	if string(wav[36:40]) != "data" {
		t.Fatalf("missing data chunk, got %q", wav[36:40])
	}

	// 44-byte header + 6 bytes of PCM (3 samples * 2 bytes)
	if len(wav) != 50 {
		t.Fatalf("len(wav) = %d, want 50", len(wav))
	}
}

func TestPCM16ToWAVRoundTripBytes(t *testing.T) {
	samples := []int16{0, 32767, -32768, 1234}
	wav, err := PCM16ToWAV(samples, 8000, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	pcm := wav[44:]
	roundTrip := pcm16ToBytes(samples)
	if len(pcm) != len(roundTrip) {
		t.Fatalf("pcm len = %d, want %d", len(pcm), len(roundTrip))
	}
	for i := range pcm {
		if pcm[i] != roundTrip[i] {
			t.Fatalf("byte %d = %d, want %d", i, pcm[i], roundTrip[i])
		}
	}
}

func TestPCM16ToBytes(t *testing.T) {
	samples := []int16{1, -1, 256, -256}
	got := pcm16ToBytes(samples)
	if len(got) != 8 {
		t.Fatalf("len = %d, want 8", len(got))
	}
	// little-endian: 1 -> [0x01, 0x00]
	if got[0] != 0x01 || got[1] != 0x00 {
		t.Errorf("first sample bytes = %d %d, want 1 0", got[0], got[1])
	}
	// -1 -> [0xFF, 0xFF]
	if got[2] != 0xFF || got[3] != 0xFF {
		t.Errorf("second sample bytes = %d %d, want 255 255", got[2], got[3])
	}
}
