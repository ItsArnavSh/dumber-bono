package entity

type Number interface {
	int16 | float32
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
