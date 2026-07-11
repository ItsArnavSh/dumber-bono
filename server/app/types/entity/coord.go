package entity

import "fmt"

type Number interface {
	int16 | float32 | uint8 | uint16
}

type Coordinates[num Number] struct {
	X, Y, Z num
}

type LatLon struct {
	Lateral, Longitudinal float32
}

type Orientation struct {
	Yaw, Pitch, Roll float32
}

type SectorWise struct {
	Sector1, Sector2, Sector3 LapTimeStamp
}

type LapTimeStamp struct {
	Minute uint8
	MS     uint16
}

func (l LapTimeStamp) String() string {
	seconds := float64(l.MS) / 1000
	return fmt.Sprintf("%d:%06.3f", l.Minute, seconds)
}

// GapString formats a delta for natural-sounding radio callouts,
// e.g. "5 tenths", "1.2 seconds", "1 minute 3 seconds".
func (l LapTimeStamp) GapString() string {
	totalMS := int(l.Minute)*60000 + int(l.MS)

	if l.Minute > 0 {
		seconds := float64(totalMS%60000) / 1000
		return fmt.Sprintf("%d minute %.1f seconds", l.Minute, seconds)
	}

	if totalMS < 1000 {
		tenths := totalMS / 100
		if tenths == 0 {
			return "less than a tenth"
		}
		if tenths == 1 {
			return "1 tenth"
		}
		return fmt.Sprintf("%d tenths", tenths)
	}

	seconds := float64(totalMS) / 1000
	return fmt.Sprintf("%.1f seconds", seconds)
}

type Tyres[num Number] struct {
	RL, RR, FL, FR num
}
