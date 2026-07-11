package entity

import "time"

type PilotPressurePhysicalFactors struct {
	// Turn Indicators
	GForceLat   float32
	GFroceLon   float32
	Steer       float32
	Brake       float32
	LapDistance float32
}

type RadioMessage struct {
	Priority int
	Message  string
	Expiry   time.Time
}

func (r RadioMessage) GetExpiry() time.Time {
	return r.Expiry
}
