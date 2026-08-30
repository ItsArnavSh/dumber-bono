package mappers

import (
	"dubmer-bono/app/api/udp/parsers"
	"testing"
)

func TestMapToLapData(t *testing.T) {
	src := &parsers.LapData{
		LastLapTimeInMS:              95000,
		CurrentLapTimeInMS:           20000,
		Sector1TimeMSPart:            25000,
		Sector1TimeMinutesPart:       0,
		Sector2TimeMSPart:            45000,
		Sector2TimeMinutesPart:       0,
		DeltaToCarInFrontMSPart:      500,
		DeltaToCarInFrontMinutesPart: 0,
		DeltaToRaceLeaderMSPart:      5000,
		DeltaToRaceLeaderMinutesPart: 0,
		LapDistance:                  3000.5,
		TotalDistance:                60000,
		CarPosition:                  3,
		CurrentLapNum:                12,
		PitStatus:                    1,
		CurrentLapInvalid:            0,
		TotalWarnings:                2,
		SpeedTrapFastestSpeed:        320.5,
	}

	got, err := MapToLapData(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.LastLapTime != 95000 {
		t.Errorf("LastLapTime = %d, want 95000", got.LastLapTime)
	}
	if got.CarPosition != 3 {
		t.Errorf("CarPosition = %d, want 3", got.CarPosition)
	}
	if got.CurrentLapNum != 12 {
		t.Errorf("CurrentLapNum = %d, want 12", got.CurrentLapNum)
	}
	if got.DeltaToFront.MS != 500 {
		t.Errorf("DeltaToFront.MS = %d, want 500", got.DeltaToFront.MS)
	}
	if got.DeltaToRaceLeader.MS != 5000 {
		t.Errorf("DeltaToRaceLeader.MS = %d, want 5000", got.DeltaToRaceLeader.MS)
	}
	if got.CurrentLapInvalid {
		t.Error("CurrentLapInvalid = true, want false")
	}
	if got.SectorTimes.Sector1.MS != 25000 {
		t.Errorf("Sector1.MS = %d, want 25000", got.SectorTimes.Sector1.MS)
	}
}

func TestMapToLapDataInvalidLap(t *testing.T) {
	src := &parsers.LapData{CurrentLapInvalid: 1}
	got, err := MapToLapData(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.CurrentLapInvalid {
		t.Error("CurrentLapInvalid = false, want true")
	}
}

func TestMapToLapDataNil(t *testing.T) {
	if _, err := MapToLapData(nil); err == nil {
		t.Fatal("expected error for nil input")
	}
}

func TestMapToLapDataPacket(t *testing.T) {
	src := &parsers.PacketLapData{}
	src.Header.PacketID = 2
	src.LapData[0].CarPosition = 1
	src.LapData[0].CurrentLapNum = 5

	got, err := MapToLapDataPacket(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Header.PacketID != 2 {
		t.Errorf("Header.PacketID = %d, want 2", got.Header.PacketID)
	}
	if got.LapData[0].CarPosition != 1 {
		t.Errorf("LapData[0].CarPosition = %d, want 1", got.LapData[0].CarPosition)
	}
	if got.LapData[0].CurrentLapNum != 5 {
		t.Errorf("LapData[0].CurrentLapNum = %d, want 5", got.LapData[0].CurrentLapNum)
	}
}

func TestMapToLapDataPacketNil(t *testing.T) {
	if _, err := MapToLapDataPacket(nil); err == nil {
		t.Fatal("expected error for nil input")
	}
}
