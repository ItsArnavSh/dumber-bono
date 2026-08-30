package parsers

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestParsePacketLapData(t *testing.T) {
	var packet PacketLapData
	packet.Header.PacketID = 2
	packet.Header.SessionUID = 12345
	packet.LapData[3].CarPosition = 5
	packet.LapData[3].CurrentLapNum = 12
	packet.TimeTrialPBCarIdx = 255
	packet.TimeTrialRivalCarIdx = 0

	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, &packet); err != nil {
		t.Fatalf("encode: %v", err)
	}

	got, err := ParsePacket[PacketLapData](buf.Bytes())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Header.PacketID != 2 {
		t.Errorf("Header.PacketID = %d, want 2", got.Header.PacketID)
	}
	if got.Header.SessionUID != 12345 {
		t.Errorf("Header.SessionUID = %d, want 12345", got.Header.SessionUID)
	}
	if got.LapData[3].CarPosition != 5 {
		t.Errorf("LapData[3].CarPosition = %d, want 5", got.LapData[3].CarPosition)
	}
	if got.LapData[3].CurrentLapNum != 12 {
		t.Errorf("LapData[3].CurrentLapNum = %d, want 12", got.LapData[3].CurrentLapNum)
	}
}

func TestParsePacketTooShort(t *testing.T) {
	if _, err := ParsePacket[PacketLapData](make([]byte, 10)); err == nil {
		t.Fatal("expected error for too-short payload")
	}
}

func TestParsePacketParticipants(t *testing.T) {
	var packet PacketParticipantsData
	packet.NumActiveCars = 20
	packet.Participants[0].DriverId = 10
	packet.Participants[0].RaceNumber = 44
	copy(packet.Participants[0].Name[:], "Lewis")

	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, &packet); err != nil {
		t.Fatalf("encode: %v", err)
	}

	got, err := ParsePacket[PacketParticipantsData](buf.Bytes())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.NumActiveCars != 20 {
		t.Errorf("NumActiveCars = %d, want 20", got.NumActiveCars)
	}
	if got.Participants[0].DriverId != 10 {
		t.Errorf("Participants[0].DriverId = %d, want 10", got.Participants[0].DriverId)
	}
	if got.Participants[0].RaceNumber != 44 {
		t.Errorf("Participants[0].RaceNumber = %d, want 44", got.Participants[0].RaceNumber)
	}
	if string(got.Participants[0].Name[:5]) != "Lewis" {
		t.Errorf("Participants[0].Name = %v, want Lewis", got.Participants[0].Name)
	}
}
