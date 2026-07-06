package mappers

import (
	"dubmer-bono/app/api/udp/parsers"
	"dubmer-bono/app/types/entity"
	telentity "dubmer-bono/app/types/entity/tel-entity"
	"errors"
)

func MapToMotion(data *parsers.CarMotionData) (telentity.CarMotion, error) {
	if data == nil {
		return telentity.CarMotion{}, errors.New("mappers: nil CarMotionData")
	}

	return telentity.CarMotion{
		WorldPosition:   entity.Coordinates[float32]{X: data.WorldPositionX, Y: data.WorldPositionY, Z: data.WorldPositionZ},
		WorldVelocity:   entity.Coordinates[float32]{X: data.WorldVelocityX, Y: data.WorldVelocityY, Z: data.WorldVelocityZ},
		WorldForwardDir: entity.Coordinates[int16]{X: data.WorldForwardDirX, Y: data.WorldForwardDirY, Z: data.WorldForwardDirZ},
		WorldRightDir:   entity.Coordinates[int16]{X: data.WorldRightDirX, Y: data.WorldRightDirY, Z: data.WorldRightDirZ},
		GForce:          entity.LatLon{Lateral: data.GForceLateral, Longitudinal: data.GForceLongitudinal},
		Orientation:     entity.Orientation{Yaw: data.Yaw, Pitch: data.Pitch, Roll: data.Roll},
	}, nil
}
func MaptoMotionPacket(data *parsers.PacketMotionData) (telentity.MotionPacket, error) {
	if data == nil {
		return telentity.MotionPacket{}, errors.New("mappers: nil PacketMotionData")
	}

	header, err := MapToHeader(&data.Header)
	if err != nil {
		return telentity.MotionPacket{}, err
	}

	var cars [22]telentity.CarMotion
	for i := range data.CarMotionData {
		car, err := MapToMotion(&data.CarMotionData[i])
		if err != nil {
			return telentity.MotionPacket{}, err
		}
		cars[i] = car
	}

	return telentity.MotionPacket{
		Header: header,
		Cars:   cars,
	}, nil
}
