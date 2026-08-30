package entity

import (
	"strings"
	"testing"
)

func TestLapTimeStampString(t *testing.T) {
	tests := []struct {
		name string
		ts   LapTimeStamp
		want string
	}{
		{"zero", LapTimeStamp{Minute: 0, MS: 0}, "0:00.000"},
		{"under a second", LapTimeStamp{Minute: 0, MS: 500}, "0:00.500"},
		{"over a minute", LapTimeStamp{Minute: 1, MS: 23456}, "1:23.456"},
		{"two minutes", LapTimeStamp{Minute: 2, MS: 0}, "2:00.000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ts.String(); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLapTimeStampGapString(t *testing.T) {
	tests := []struct {
		name string
		ts   LapTimeStamp
		want string
	}{
		{"zero", LapTimeStamp{Minute: 0, MS: 0}, "less than a tenth"},
		{"one tenth", LapTimeStamp{Minute: 0, MS: 100}, "1 tenth"},
		{"five tenths", LapTimeStamp{Minute: 0, MS: 500}, "5 tenths"},
		{"just under a second", LapTimeStamp{Minute: 0, MS: 900}, "9 tenths"},
		{"one second", LapTimeStamp{Minute: 0, MS: 1000}, "1.0 seconds"},
		{"twelve tenths rounds to seconds", LapTimeStamp{Minute: 0, MS: 1200}, "1.2 seconds"},
		{"over a minute", LapTimeStamp{Minute: 1, MS: 3000}, "1 minute 3.0 seconds"},
		{"two minutes zero", LapTimeStamp{Minute: 2, MS: 0}, "2 minute 0.0 seconds"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ts.GapString()
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLapTimeStampGapStringHasMinute(t *testing.T) {
	ts := LapTimeStamp{Minute: 3, MS: 999}
	if !strings.Contains(ts.GapString(), "3 minute") {
		t.Fatalf("expected minute unit in %q", ts.GapString())
	}
}

func TestCoordinatesZero(t *testing.T) {
	var c Coordinates[float32]
	if c.X != 0 || c.Y != 0 || c.Z != 0 {
		t.Fatal("expected zero-valued coordinates")
	}
}
