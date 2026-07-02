package entity

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

type Tyres[num Number] struct {
	RL, RR, FL, FR num
}
