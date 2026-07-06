package mappers

import (
	"errors"

	"dubmer-bono/app/api/udp/parsers"
	"dubmer-bono/app/types/entity"
	telentity "dubmer-bono/app/types/entity/tel-entity"
)

func MapToCarDamage(data *parsers.CarDamageData) (telentity.CarDamageData, error) {
	if data == nil {
		return telentity.CarDamageData{}, errors.New("mappers: nil CarDamageData")
	}

	return telentity.CarDamageData{
		TyresWear:            entity.Tyres[float32]{RL: data.TyresWear[0], RR: data.TyresWear[1], FL: data.TyresWear[2], FR: data.TyresWear[3]},
		TyresDamage:          entity.Tyres[uint8]{RL: data.TyresDamage[0], RR: data.TyresDamage[1], FL: data.TyresDamage[2], FR: data.TyresDamage[3]},
		BrakesDamage:         entity.Tyres[uint8]{RL: data.BrakesDamage[0], RR: data.BrakesDamage[1], FL: data.BrakesDamage[2], FR: data.BrakesDamage[3]},
		TyreBlisters:         entity.Tyres[uint8]{RL: data.TyreBlisters[0], RR: data.TyreBlisters[1], FL: data.TyreBlisters[2], FR: data.TyreBlisters[3]},
		FrontLeftWingDamage:  data.FrontLeftWingDamage,
		FrontRightWingDamage: data.FrontRightWingDamage,
		RearWingDamage:       data.RearWingDamage,
		FloorDamage:          data.FloorDamage,
		DiffuserDamage:       data.DiffuserDamage,
		SidepodDamage:        data.SidepodDamage,
		DRSFault:             data.DrsFault != 0,
		ERSFault:             data.ErsFault != 0,
		GearBoxDamage:        data.GearBoxDamage,
		EngineDamage:         data.EngineDamage,
		EngineMGUHWear:       data.EngineMGUHWear,
		EngineESWear:         data.EngineESWear,
		EngineCEWear:         data.EngineCEWear,
		EngineICEWear:        data.EngineICEWear,
		EngineMGUKWear:       data.EngineMGUKWear,
		EngineTCWear:         data.EngineTCWear,
		EngineBlown:          data.EngineBlown != 0,
		EngineSeized:         data.EngineSeized != 0,
	}, nil
}

func MapToCarDamagePacket(data *parsers.PacketCarDamageData) (telentity.PacketCarDamageData, error) {
	if data == nil {
		return telentity.PacketCarDamageData{}, errors.New("mappers: nil PacketCarDamageData")
	}

	header, err := MapToHeader(&data.Header)
	if err != nil {
		return telentity.PacketCarDamageData{}, err
	}

	var cars [22]telentity.CarDamageData
	for i := range data.CarDamageData {
		car, err := MapToCarDamage(&data.CarDamageData[i])
		if err != nil {
			return telentity.PacketCarDamageData{}, err
		}
		cars[i] = car
	}

	return telentity.PacketCarDamageData{
		Header:        header,
		CarDamageData: cars,
	}, nil
}
