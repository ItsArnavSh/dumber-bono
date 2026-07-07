package mappers

import (
	"errors"

	"dubmer-bono/app/api/udp/parsers"
	telentity "dubmer-bono/app/types/entity/tel-entity"
)

func MapToCarSetup(data *parsers.CarSetupData) (telentity.CarSetupData, error) {
	if data == nil {
		return telentity.CarSetupData{}, errors.New("mappers: nil CarSetupData")
	}

	return telentity.CarSetupData{
		FrontWing:              data.FrontWing,
		RearWing:               data.RearWing,
		OnThrottle:             data.OnThrottle,
		OffThrottle:            data.OffThrottle,
		FrontCamber:            data.FrontCamber,
		RearCamber:             data.RearCamber,
		FrontToe:               data.FrontToe,
		RearToe:                data.RearToe,
		FrontSuspension:        data.FrontSuspension,
		RearSuspension:         data.RearSuspension,
		FrontAntiRollBar:       data.FrontAntiRollBar,
		RearAntiRollBar:        data.RearAntiRollBar,
		FrontSuspensionHeight:  data.FrontSuspensionHeight,
		RearSuspensionHeight:   data.RearSuspensionHeight,
		BrakePressure:          data.BrakePressure,
		BrakeBias:              data.BrakeBias,
		EngineBraking:          data.EngineBraking,
		RearLeftTyrePressure:   data.RearLeftTyrePressure,
		RearRightTyrePressure:  data.RearRightTyrePressure,
		FrontLeftTyrePressure:  data.FrontLeftTyrePressure,
		FrontRightTyrePressure: data.FrontRightTyrePressure,
		Ballast:                data.Ballast,
		FuelLoad:               data.FuelLoad,
	}, nil
}

func MapToCarSetupPacket(data *parsers.PacketCarSetupData) (telentity.PacketCarSetupData, error) {
	if data == nil {
		return telentity.PacketCarSetupData{}, errors.New("mappers: nil PacketCarSetupData")
	}

	header, err := MapToHeader(&data.Header)
	if err != nil {
		return telentity.PacketCarSetupData{}, err
	}

	var setups [22]telentity.CarSetupData
	for i := range data.CarSetups {
		setup, err := MapToCarSetup(&data.CarSetups[i])
		if err != nil {
			return telentity.PacketCarSetupData{}, err
		}
		setups[i] = setup
	}

	return telentity.PacketCarSetupData{
		Header:             header,
		CarSetups:          setups,
		NextFrontWingValue: data.NextFrontWingValue,
	}, nil
}
