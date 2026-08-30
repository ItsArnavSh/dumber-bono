package mappers

import (
	"dubmer-bono/app/api/udp/parsers"
	"testing"
)

func TestMapToMotion(t *testing.T) {
	src := &parsers.CarMotionData{
		WorldPositionX:     1.5,
		WorldPositionY:     -2.5,
		WorldPositionZ:     3.0,
		WorldVelocityX:     10,
		WorldVelocityY:     20,
		WorldVelocityZ:     -5,
		WorldForwardDirX:   100,
		WorldForwardDirY:   -100,
		WorldForwardDirZ:   0,
		GForceLateral:      3.2,
		GForceLongitudinal: 1.1,
		Yaw:                0.5,
		Pitch:              -0.2,
		Roll:               0.1,
	}

	got, err := MapToMotion(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.WorldPosition.X != 1.5 {
		t.Errorf("WorldPosition.X = %v, want 1.5", got.WorldPosition.X)
	}
	if got.WorldPosition.Y != -2.5 {
		t.Errorf("WorldPosition.Y = %v, want -2.5", got.WorldPosition.Y)
	}
	if got.WorldVelocity.X != 10 {
		t.Errorf("WorldVelocity.X = %v, want 10", got.WorldVelocity.X)
	}
	if got.GForce.Lateral != 3.2 {
		t.Errorf("GForce.Lateral = %v, want 3.2", got.GForce.Lateral)
	}
	if got.Orientation.Pitch != -0.2 {
		t.Errorf("Orientation.Pitch = %v, want -0.2", got.Orientation.Pitch)
	}
}

func TestMapToMotionNil(t *testing.T) {
	if _, err := MapToMotion(nil); err == nil {
		t.Fatal("expected error for nil input")
	}
}

func TestMaptoMotionPacket(t *testing.T) {
	src := &parsers.PacketMotionData{}
	src.Header.PacketID = 0
	src.CarMotionData[3].WorldPositionX = 99.5

	got, err := MaptoMotionPacket(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Header.PacketID != 0 {
		t.Errorf("Header.PacketID = %d, want 0", got.Header.PacketID)
	}
	if got.Cars[3].WorldPosition.X != 99.5 {
		t.Errorf("Cars[3].WorldPosition.X = %v, want 99.5", got.Cars[3].WorldPosition.X)
	}
}

func TestMaptoMotionPacketNil(t *testing.T) {
	if _, err := MaptoMotionPacket(nil); err == nil {
		t.Fatal("expected error for nil input")
	}
}
