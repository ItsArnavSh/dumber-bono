package mappers

import (
	"errors"

	"dubmer-bono/app/api/udp/parsers"
	"dubmer-bono/app/types/entity/consts"
	telentity "dubmer-bono/app/types/entity/tel-entity"
)

func MapToTyreSet(data *parsers.TyreSetData) (telentity.TyreSetData, error) {
	if data == nil {
		return telentity.TyreSetData{}, errors.New("mappers: nil TyreSetData")
	}

	return telentity.TyreSetData{
		ActualTyreCompound: consts.ActualTyreCompound[int16(data.ActualTyreCompound)],
		VisualTyreCompound: consts.VisualTyreCompound[int16(data.VisualTyreCompound)],
		Wear:               data.Wear,
		Available:          data.Available != 0,
		RecommendedSession: data.RecommendedSession,
		LifeSpan:           data.LifeSpan,
		UsableLife:         data.UsableLife,
		LapDeltaTime:       data.LapDeltaTime,
		Fitted:             data.Fitted != 0,
	}, nil
}

func MapToTyreSetPacket(data *parsers.PacketTyreSetsData) (telentity.TyreSetPacket, error) {
	if data == nil {
		return telentity.TyreSetPacket{}, errors.New("mappers: nil PacketTyreSetsData")
	}

	header, err := MapToHeader(&data.Header)
	if err != nil {
		return telentity.TyreSetPacket{}, err
	}

	var sets [20]telentity.TyreSetData
	for i := range data.TyreSetData {
		set, err := MapToTyreSet(&data.TyreSetData[i])
		if err != nil {
			return telentity.TyreSetPacket{}, err
		}
		sets[i] = set
	}

	return telentity.TyreSetPacket{
		Header:     header,
		Carid:      data.CarIdx,
		Tyredata:   sets,
		Fittedtyre: data.FittedIdx,
	}, nil
}
