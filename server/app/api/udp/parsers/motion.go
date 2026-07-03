package parsers

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

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

// ParseMotionPacket decodes a raw UDP payload into a PacketMotionData.
func ParseMotionPacket(payload []byte) (*PacketMotionData, error) {
	var packet PacketMotionData

	reader := bytes.NewReader(payload)
	if err := binary.Read(reader, binary.LittleEndian, &packet); err != nil {
		return nil, fmt.Errorf("decode motion packet: %w", err)
	}

	return &packet, nil
}
