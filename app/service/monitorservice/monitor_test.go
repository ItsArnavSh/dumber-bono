package monitor

import (
	"dubmer-bono/app/types/entity"
	"testing"
)

func TestIsNearPosition(t *testing.T) {
	tests := []struct {
		name     string
		mypos    int
		pos      int
		expected bool
	}{
		{"same position", 5, 5, true},
		{"one ahead", 5, 4, true},
		{"one behind", 5, 6, true},
		{"two ahead", 5, 3, true},
		{"two behind", 5, 7, true},
		{"three ahead", 5, 2, false},
		{"three behind", 5, 8, false},
		{"edge zero", 1, 3, true},
		{"edge upper", 20, 18, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNearPosition(tt.mypos, tt.pos); got != tt.expected {
				t.Fatalf("isNearPosition(%d, %d) = %v, want %v", tt.mypos, tt.pos, got, tt.expected)
			}
		})
	}
}

func TestAlterConfidence(t *testing.T) {
	tests := []struct {
		name     string
		pressure entity.PilotPressurePhysicalFactors
		want     int
	}{
		{"no pressure", entity.PilotPressurePhysicalFactors{}, 0},
		{"full braking", entity.PilotPressurePhysicalFactors{Brake: 1.0}, 3},
		{"moderate steer", entity.PilotPressurePhysicalFactors{Steer: 0.4}, 1},
		{"max everything", entity.PilotPressurePhysicalFactors{GForceLat: 5.0, GFroceLon: 5.0, Steer: 1.0, Brake: 1.0}, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{}
			s.AlterConfidence(tt.pressure)
			if s.driver_pressure != tt.want {
				t.Fatalf("driver_pressure = %d, want %d", s.driver_pressure, tt.want)
			}
		})
	}
}

func TestAlterConfidenceBounds(t *testing.T) {
	s := &Service{}
	// Excessive values must still clamp to 0..5
	s.AlterConfidence(entity.PilotPressurePhysicalFactors{
		GForceLat: 50.0, GFroceLon: 50.0, Steer: 10.0, Brake: 2.0,
	})
	if s.driver_pressure < 0 || s.driver_pressure > 5 {
		t.Fatalf("driver_pressure = %d out of range 0..5", s.driver_pressure)
	}
}
