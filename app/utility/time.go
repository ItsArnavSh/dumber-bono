package utility

import (
	"fmt"
	"time"
)

// CheckFramesFresh takes any number of event timestamps and the expected
// max spread between them. It returns the midpoint time across all events,
// and an error if the spread between earliest and latest exceeds tolerance.
func CheckFramesFresh(expected time.Duration, times ...time.Time) (avg time.Time, err error) {
	if len(times) == 0 {
		return time.Time{}, fmt.Errorf("no timestamps provided")
	}

	earliest, latest := times[0], times[0]
	for _, t := range times[1:] {
		if t.Before(earliest) {
			earliest = t
		}
		if t.After(latest) {
			latest = t
		}
	}

	spread := latest.Sub(earliest)
	avg = earliest.Add(spread / 2)

	const tolerance = 1.5
	if spread > time.Duration(float64(expected)*tolerance) {
		return avg, fmt.Errorf("spread %v exceeds expected %v", spread, expected)
	}
	return avg, nil
}
