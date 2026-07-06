package mappers

import (
	"errors"

	"dubmer-bono/app/api/udp/parsers"
	telentity "dubmer-bono/app/types/entity/tel-entity"
)

func MapToCarTelemetry(data *parsers.CarTelemetryData) (telentity.CarTelemetryData, error) {
	if data == nil {
		return telentity.CarTelemetryData{}, errors.New("mappers: nil CarTelemetryData")
	}

	return telentity.CarTelemetryData{
		Speed:                   data.Speed,
		Throttle:                data.Throttle,
		Steer:                   data.Steer,
		Brake:                   data.Brake,
		Clutch:                  data.Clutch,
		Gear:                    data.Gear,
		EngineRPM:               data.EngineRPM,
		DRS:                     data.Drs != 0,
		BrakesTemperature:       data.BrakesTemperature,
		TyresSurfaceTemperature: data.TyresSurfaceTemperature,
		TyresInnerTemperature:   data.TyresInnerTemperature,
		EngineTemperature:       data.EngineTemperature,
		TyresPressure:           data.TyresPressure,
		SurfaceType:             data.SurfaceType,
	}, nil
}

func MapToCarTelemetryPacket(data *parsers.PacketCarTelemetryData) (telentity.PacketCarTelemetryData, error) {
	if data == nil {
		return telentity.PacketCarTelemetryData{}, errors.New("mappers: nil PacketCarTelemetryData")
	}

	header, err := MapToHeader(&data.Header)
	if err != nil {
		return telentity.PacketCarTelemetryData{}, err
	}

	var cars [22]telentity.CarTelemetryData
	for i := range data.CarTelemetryData {
		car, err := MapToCarTelemetry(&data.CarTelemetryData[i])
		if err != nil {
			return telentity.PacketCarTelemetryData{}, err
		}
		cars[i] = car
	}

	return telentity.PacketCarTelemetryData{
		Header:                       header,
		CarTelemetryData:             cars,
		MFDPanelIndex:                data.MfdPanelIndex,
		MFDPanelIndexSecondaryPlayer: data.MfdPanelIndexSecondaryPlayer,
		SuggestedGear:                data.SuggestedGear,
	}, nil
}
