package entity

import "time"

type TelemetryFrame struct {
	SessionID uint32
	CarNo     uint8
	FrameTime time.Time

	Speed     float32
	Throttle  float32
	Steer     float32
	Brake     float32
	Clutch    float32
	Gear      int8
	EngineRPM uint16
	DRS       bool

	PosX, PosY, PosZ float32
	VelX, VelY, VelZ float32
	FwdX, FwdY, FwdZ float32

	GForceLat, GForceLon float32

	Yaw, Pitch, Roll float32
}
