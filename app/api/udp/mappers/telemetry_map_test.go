package mappers

import (
	"dubmer-bono/app/api/udp/parsers"
	"testing"
)

func TestMapToCarTelemetry(t *testing.T) {
	src := &parsers.CarTelemetryData{
		Speed:             200,
		Throttle:          85,
		Steer:             0.5,
		Brake:             30,
		Clutch:            0,
		Gear:              7,
		EngineRPM:         11500,
		Drs:               1,
		EngineTemperature: 100,
		SurfaceType:       [4]uint8{3, 3, 3, 3},
	}

	got, err := MapToCarTelemetry(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Speed != 200 {
		t.Errorf("Speed = %d, want 200", got.Speed)
	}
	if got.DRS != true {
		t.Errorf("DRS = %v, want true", got.DRS)
	}
	if got.Gear != 7 {
		t.Errorf("Gear = %d, want 7", got.Gear)
	}
	if got.SurfaceType[0] != 3 {
		t.Errorf("SurfaceType = %v, want [3 3 3 3]", got.SurfaceType)
	}
}

func TestMapToCarTelemetryDRSOff(t *testing.T) {
	src := &parsers.CarTelemetryData{Drs: 0}
	got, err := MapToCarTelemetry(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.DRS {
		t.Error("DRS = true, want false")
	}
}

func TestMapToCarTelemetryNil(t *testing.T) {
	if _, err := MapToCarTelemetry(nil); err == nil {
		t.Fatal("expected error for nil input")
	}
}

func TestMapToCarTelemetryPacket(t *testing.T) {
	src := &parsers.PacketCarTelemetryData{}
	src.Header.PacketID = 6
	src.Header.SessionUID = 77
	src.CarTelemetryData[1].Speed = 250
	src.MfdPanelIndex = 1
	src.MfdPanelIndexSecondaryPlayer = 255
	src.SuggestedGear = 7

	got, err := MapToCarTelemetryPacket(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Header.PacketID != 6 {
		t.Errorf("Header.PacketID = %d, want 6", got.Header.PacketID)
	}
	if got.CarTelemetryData[1].Speed != 250 {
		t.Errorf("CarTelemetryData[1].Speed = %d, want 250", got.CarTelemetryData[1].Speed)
	}
	if got.SuggestedGear != 7 {
		t.Errorf("SuggestedGear = %d, want 7", got.SuggestedGear)
	}
}

func TestMapToCarTelemetryPacketNil(t *testing.T) {
	if _, err := MapToCarTelemetryPacket(nil); err == nil {
		t.Fatal("expected error for nil input")
	}
}
