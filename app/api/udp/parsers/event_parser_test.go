package parsers

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func encodeEventPacket(t *testing.T, header PacketHeader, code [4]byte, details any) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, &header); err != nil {
		t.Fatalf("encode header: %v", err)
	}
	if err := binary.Write(&buf, binary.LittleEndian, code); err != nil {
		t.Fatalf("encode code: %v", err)
	}
	if details != nil {
		if err := binary.Write(&buf, binary.LittleEndian, details); err != nil {
			t.Fatalf("encode details: %v", err)
		}
	}
	return buf.Bytes()
}

func TestParseEventPacketFastestLap(t *testing.T) {
	payload := encodeEventPacket(t, PacketHeader{PacketID: 3}, [4]byte{'F', 'T', 'L', 'P'}, FastestLap{VehicleIdx: 7, LapTime: 91.5})

	packet, err := ParseEventPacket(payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if packet.Header.PacketID != 3 {
		t.Errorf("PacketID = %d, want 3", packet.Header.PacketID)
	}
	if got := packet.EventStringCode; got != [4]uint8{'F', 'T', 'L', 'P'} {
		t.Errorf("EventStringCode = %v, want %v", got, [4]uint8{'F', 'T', 'L', 'P'})
	}
	if packet.EventDetails.FastestLap.VehicleIdx != 7 {
		t.Errorf("FastestLap.VehicleIdx = %d, want 7", packet.EventDetails.FastestLap.VehicleIdx)
	}
	if packet.EventDetails.FastestLap.LapTime != 91.5 {
		t.Errorf("FastestLap.LapTime = %v, want 91.5", packet.EventDetails.FastestLap.LapTime)
	}
}

func TestParseEventPacketRetirement(t *testing.T) {
	payload := encodeEventPacket(t, PacketHeader{}, [4]byte{'R', 'T', 'M', 'T'}, Retirement{VehicleIdx: 4, Reason: 1})

	packet, err := ParseEventPacket(payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if packet.EventDetails.Retirement.VehicleIdx != 4 {
		t.Errorf("Retirement.VehicleIdx = %d, want 4", packet.EventDetails.Retirement.VehicleIdx)
	}
	if packet.EventDetails.Retirement.Reason != 1 {
		t.Errorf("Retirement.Reason = %d, want 1", packet.EventDetails.Retirement.Reason)
	}
}

func TestParseEventPacketPenalty(t *testing.T) {
	details := Penalty{PenaltyType: 3, InfringementType: 1, VehicleIdx: 2, OtherVehicleIdx: 5, Time: 5, LapNum: 12, PlacesGained: 1}
	payload := encodeEventPacket(t, PacketHeader{}, [4]byte{'P', 'E', 'N', 'A'}, details)

	packet, err := ParseEventPacket(payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := packet.EventDetails.Penalty
	if got != details {
		t.Errorf("Penalty = %+v, want %+v", got, details)
	}
}

func TestParseEventPacketNoPayload(t *testing.T) {
	payload := encodeEventPacket(t, PacketHeader{}, [4]byte{'S', 'S', 'T', 'A'}, nil)

	packet, err := ParseEventPacket(payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := string(packet.EventStringCode[:]); got != "SSTA" {
		t.Errorf("EventStringCode = %q, want %q", got, "SSTA")
	}
}

func TestParseEventPacketUnknownCode(t *testing.T) {
	payload := encodeEventPacket(t, PacketHeader{}, [4]byte{'X', 'X', 'X', 'X'}, nil)

	if _, err := ParseEventPacket(payload); err == nil {
		t.Fatal("expected error for unknown event code")
	}
}

func TestParseEventPacketTruncated(t *testing.T) {
	payload := encodeEventPacket(t, PacketHeader{}, [4]byte{'F', 'T', 'L', 'P'}, FastestLap{VehicleIdx: 1, LapTime: 90})
	if _, err := ParseEventPacket(payload[:len(payload)-1]); err == nil {
		t.Fatal("expected error for truncated payload")
	}
}

func TestBytesToEventCode(t *testing.T) {
	tests := []struct {
		name string
		in   [4]uint8
		want string
	}{
		{"full", [4]uint8{'F', 'T', 'L', 'P'}, "FTLP"},
		{"nul-terminated", [4]uint8{'S', 'C', 'A', 'R'}, "SCAR"},
		{"all zero", [4]uint8{0, 0, 0, 0}, ""},
		{"partial", [4]uint8{'S', 'S', 'T', 'A'}, "SSTA"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bytesToEventCode(tt.in); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
