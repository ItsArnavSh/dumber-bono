package mappers

import (
	"dubmer-bono/app/api/udp/parsers"
	telentity "dubmer-bono/app/types/entity/tel-entity"
	"testing"
)

func TestMapToEventDataFastestLap(t *testing.T) {
	src := &parsers.PacketEventData{
		EventStringCode: [4]uint8{'F', 'T', 'L', 'P'},
		EventDetails:    parsers.EventDataDetails{FastestLap: parsers.FastestLap{VehicleIdx: 2, LapTime: 92.4}},
	}

	got, err := MapToEventData(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.EventStringCode != "FTLP" {
		t.Errorf("EventStringCode = %q, want FTLP", got.EventStringCode)
	}
	d, ok := got.Details.(telentity.FastestLap)
	if !ok {
		t.Fatalf("Details type = %T, want telentity.FastestLap", got.Details)
	}
	if d.VehicleIdx != 2 || d.LapTime != 92.4 {
		t.Errorf("FastestLap = %+v, want VehicleIdx=2 LapTime=92.4", d)
	}
}

func TestMapToEventDataSpeedTrap(t *testing.T) {
	src := &parsers.PacketEventData{
		EventStringCode: [4]uint8{'S', 'P', 'T', 'P'},
		EventDetails: parsers.EventDataDetails{
			SpeedTrap: parsers.SpeedTrap{
				VehicleIdx:                 5,
				Speed:                      335.2,
				IsOverallFastestInSession:  1,
				FastestVehicleIdxInSession: 5,
				FastestSpeedInSession:      335.2,
			},
		},
	}

	got, err := MapToEventData(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	d, ok := got.Details.(telentity.SpeedTrap)
	if !ok {
		t.Fatalf("Details type = %T, want telentity.SpeedTrap", got.Details)
	}
	if !d.IsOverallFastestInSession {
		t.Error("IsOverallFastestInSession = false, want true")
	}
	if d.Speed != 335.2 {
		t.Errorf("Speed = %v, want 335.2", d.Speed)
	}
}

func TestMapToEventDataNil(t *testing.T) {
	if _, err := MapToEventData(nil); err == nil {
		t.Fatal("expected error for nil input")
	}
}

func TestMapToEventDataUnknownCode(t *testing.T) {
	src := &parsers.PacketEventData{EventStringCode: [4]uint8{'Z', 'Z', 'Z', 'Z'}}
	if _, err := MapToEventData(src); err == nil {
		t.Fatal("expected error for unknown event code")
	}
}

func TestMapToEventDataNoPayload(t *testing.T) {
	src := &parsers.PacketEventData{EventStringCode: [4]uint8{'S', 'S', 'T', 'A'}}
	got, err := MapToEventData(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Details != nil {
		t.Errorf("Details = %v, want nil", got.Details)
	}
}
