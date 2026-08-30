package mappers

import (
	"dubmer-bono/app/api/udp/parsers"
	"testing"
)

func TestMapToHeader(t *testing.T) {
	src := &parsers.PacketHeader{
		PacketFormat:            2025,
		GameYear:                25,
		GameMajorVersion:        1,
		GameMinorVersion:        0,
		PacketVersion:           1,
		PacketID:                4,
		SessionUID:              999,
		SessionTime:             45.0,
		FrameIdentifier:         100,
		OverallFrameIdentifier:  200,
		PlayerCarIndex:          3,
		SecondaryPlayerCarIndex: 255,
	}

	got, err := MapToHeader(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.PacketID != 4 {
		t.Errorf("PacketID = %d, want 4", got.PacketID)
	}
	if got.SessionUID != 999 {
		t.Errorf("SessionUID = %d, want 999", got.SessionUID)
	}
	if got.PlayerCarIndex != 3 {
		t.Errorf("PlayerCarIndex = %d, want 3", got.PlayerCarIndex)
	}
}

func TestMapToHeaderNil(t *testing.T) {
	if _, err := MapToHeader(nil); err == nil {
		t.Fatal("expected error for nil header")
	}
}
