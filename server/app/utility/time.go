package utility

import (
	"fmt"
	"time"
)

// CheckFrameGap takes two event timestamps (in either order) and the
// expected duration between them. It returns the midpoint time between
// the two events, and an error if the gap deviates too far from expected.
func CheckFrameGap(event1, event2 time.Time, expected time.Duration) (avg time.Time, err error) {
	gap := event2.Sub(event1)
	if gap < 0 {
		gap = -gap
	}

	// midpoint = earlier time + half the gap
	earlier := event1
	if event2.Before(event1) {
		earlier = event2
	}
	avg = earlier.Add(gap / 2)

	const tolerance = 1.5 // allow up to 1.5x expected before flagging
	if gap > time.Duration(float64(expected)*tolerance) {
		return avg, fmt.Errorf("gap %v exceeds expected %v", gap, expected)
	}
	return avg, nil
}
