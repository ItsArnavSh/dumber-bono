package parsers

// CarMotionData mirrors the C struct CarMotionData field-for-field.
type CarMotionData struct {
	WorldPositionX     float32
	WorldPositionY     float32
	WorldPositionZ     float32
	WorldVelocityX     float32
	WorldVelocityY     float32
	WorldVelocityZ     float32
	WorldForwardDirX   int16
	WorldForwardDirY   int16
	WorldForwardDirZ   int16
	WorldRightDirX     int16
	WorldRightDirY     int16
	WorldRightDirZ     int16
	GForceLateral      float32
	GForceLongitudinal float32
	GForceVertical     float32
	Yaw                float32
	Pitch              float32
	Roll               float32
}

// PacketMotionData mirrors PacketMotionData, motion data for all 22 cars.
type PacketMotionData struct {
	Header        PacketHeader
	CarMotionData [22]CarMotionData
}
