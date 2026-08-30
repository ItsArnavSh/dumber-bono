package utility

import (
	"testing"
	"time"
)

func TestCheckFramesFreshNoTimestamps(t *testing.T) {
	if _, err := CheckFramesFresh(time.Second); err == nil {
		t.Fatal("expected error when no timestamps provided")
	}
}

func TestCheckFramesFreshSingleTimestamp(t *testing.T) {
	ts := time.Now()
	avg, err := CheckFramesFresh(time.Second, ts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !avg.Equal(ts) {
		t.Fatalf("avg %v, want %v", avg, ts)
	}
}

func TestCheckFramesFreshWithinTolerance(t *testing.T) {
	base := time.Now()
	avg, err := CheckFramesFresh(time.Second, base, base.Add(100*time.Millisecond), base.Add(200*time.Millisecond))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := base.Add(100 * time.Millisecond)
	if !avg.Equal(want) {
		t.Fatalf("avg %v, want %v", avg, want)
	}
}

func TestCheckFramesFreshExceedsTolerance(t *testing.T) {
	base := time.Now()
	if _, err := CheckFramesFresh(time.Second, base, base.Add(5*time.Second)); err == nil {
		t.Fatal("expected error when spread exceeds tolerance")
	}
}

func TestCheckFramesFreshUnorderedInput(t *testing.T) {
	base := time.Now()
	late := base.Add(200 * time.Millisecond)
	early := base.Add(-200 * time.Millisecond)

	avg, err := CheckFramesFresh(time.Second, late, early)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := base
	if !avg.Equal(want) {
		t.Fatalf("avg %v, want %v", avg, want)
	}
}
