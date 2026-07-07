package utility

import "time"

type Throttler struct {
	Interval time.Duration
	Last     time.Time
}

func (t *Throttler) Allow() bool {
	now := time.Now()
	if now.Sub(t.Last) < t.Interval {
		return false
	}
	t.Last = now
	return true
}
