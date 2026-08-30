package parsers

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestHeaderSize(t *testing.T) {
	if HeaderSize != 29 {
		t.Fatalf("HeaderSize = %d, want 29", HeaderSize)
	}
}

func TestParseHeader(t *testing.T) {
	var h PacketHeader
	h.PacketFormat = 2025
	h.GameYear = 25
	h.PacketID = 1
	h.SessionUID = 0xDEADBEEF
	h.SessionTime = 12.5
	h.FrameIdentifier = 42
	h.PlayerCarIndex = 3

	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, &h); err != nil {
		t.Fatalf("failed to encode header: %v", err)
	}

	got, err := ParseHeader(buf.Bytes())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.PacketFormat != 2025 {
		t.Errorf("PacketFormat = %d, want 2025", got.PacketFormat)
	}
	if got.GameYear != 25 {
		t.Errorf("GameYear = %d, want 25", got.GameYear)
	}
	if got.PacketID != 1 {
		t.Errorf("PacketID = %d, want 1", got.PacketID)
	}
	if got.SessionUID != 0xDEADBEEF {
		t.Errorf("SessionUID = %d, want %d", got.SessionUID, 0xDEADBEEF)
	}
	if got.SessionTime != 12.5 {
		t.Errorf("SessionTime = %v, want 12.5", got.SessionTime)
	}
	if got.FrameIdentifier != 42 {
		t.Errorf("FrameIdentifier = %d, want 42", got.FrameIdentifier)
	}
	if got.PlayerCarIndex != 3 {
		t.Errorf("PlayerCarIndex = %d, want 3", got.PlayerCarIndex)
	}
}

func TestParseHeaderTooShort(t *testing.T) {
	if _, err := ParseHeader(make([]byte, 4)); err == nil {
		t.Fatal("expected error for too-short payload")
	}
}
