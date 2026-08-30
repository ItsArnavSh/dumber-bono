package utility

import (
	"testing"
	"time"
)

func TestThrottlerAllowsFirstCall(t *testing.T) {
	throttler := Throttler{Interval: time.Second}
	if !throttler.Allow() {
		t.Fatal("expected first call to be allowed")
	}
}

func TestThrottlerBlocksWithinInterval(t *testing.T) {
	throttler := Throttler{Interval: time.Second}
	throttler.Allow()

	if throttler.Allow() {
		t.Fatal("expected call within interval to be blocked")
	}
}

func TestThrottlerAllowsAfterInterval(t *testing.T) {
	throttler := Throttler{Interval: time.Millisecond}
	throttler.Allow()

	time.Sleep(5 * time.Millisecond)

	if !throttler.Allow() {
		t.Fatal("expected call after interval to be allowed")
	}
}

func TestThrottlerZeroIntervalAlwaysAllows(t *testing.T) {
	throttler := Throttler{Interval: 0}

	for range 5 {
		if !throttler.Allow() {
			t.Fatal("expected zero-interval throttler to always allow")
		}
	}
}
