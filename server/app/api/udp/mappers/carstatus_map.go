package mappers

import (
	"errors"

	"dubmer-bono/app/api/udp/parsers"
	"dubmer-bono/app/types/entity/consts"
	telentity "dubmer-bono/app/types/entity/tel-entity"
)

func MapToCarStatus(data *parsers.CarStatusData) (telentity.CarStatusData, error) {
	if data == nil {
		return telentity.CarStatusData{}, errors.New("mappers: nil CarStatusData")
	}

	return telentity.CarStatusData{
		FuelMix:                 consts.FuelMix[int16(data.FuelMix)],
		FrontBrakeBias:          data.FrontBrakeBias,
		PitLimiterStatus:        data.PitLimiterStatus != 0,
		FuelInTank:              data.FuelInTank,
		FuelCapacity:            data.FuelCapacity,
		FuelRemainingLaps:       data.FuelRemainingLaps,
		DRSAllowed:              data.DrsAllowed != 0,
		DRSActivationDistance:   data.DrsActivationDistance,
		ActualTyreCompound:      consts.ActualTyreCompound[int16(data.ActualTyreCompound)],
		VisualTyreCompound:      consts.VisualTyreCompound[int16(data.VisualTyreCompound)],
		TyresAgeLaps:            data.TyresAgeLaps,
		VehicleFiaFlags:         consts.VehicleFiaFlags[int16(data.VehicleFiaFlags)],
		EnginePowerICE:          data.EnginePowerICE,
		EnginePowerMGUK:         data.EnginePowerMGUK,
		ERSStoreEnergy:          data.ErsStoreEnergy,
		ERSDeployMode:           consts.ERSDeployMode[int16(data.ErsDeployMode)],
		ERSHarvestedThisLapMGUK: data.ErsHarvestedThisLapMGUK,
		ERSHarvestedThisLapMGUH: data.ErsHarvestedThisLapMGUH,
		ERSDeployedThisLap:      data.ErsDeployedThisLap,
		NetworkPaused:           data.NetworkPaused != 0,
	}, nil
}

func MapToCarStatusPacket(data *parsers.PacketCarStatusData) (telentity.PacketCarStatusData, error) {
	if data == nil {
		return telentity.PacketCarStatusData{}, errors.New("mappers: nil PacketCarStatusData")
	}

	header, err := MapToHeader(&data.Header)
	if err != nil {
		return telentity.PacketCarStatusData{}, err
	}

	var cars [22]telentity.CarStatusData
	for i := range data.CarStatusData {
		car, err := MapToCarStatus(&data.CarStatusData[i])
		if err != nil {
			return telentity.PacketCarStatusData{}, err
		}
		cars[i] = car
	}

	return telentity.PacketCarStatusData{
		Header:        header,
		CarStatusData: cars,
	}, nil
}
