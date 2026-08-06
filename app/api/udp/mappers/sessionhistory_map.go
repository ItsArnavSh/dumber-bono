package mappers

import (
	"errors"

	"dubmer-bono/app/api/udp/parsers"
	"dubmer-bono/app/types/entity/consts"
	telentity "dubmer-bono/app/types/entity/tel-entity"
)

func MapToLapHistory(data *parsers.LapHistoryData) (telentity.LapHistoryData, error) {
	if data == nil {
		return telentity.LapHistoryData{}, errors.New("mappers: nil LapHistoryData")
	}

	return telentity.LapHistoryData{
		LapTimeInMS:            data.LapTimeInMS,
		Sector1TimeMSPart:      data.Sector1TimeMSPart,
		Sector1TimeMinutesPart: data.Sector1TimeMinutesPart,
		Sector2TimeMSPart:      data.Sector2TimeMSPart,
		Sector2TimeMinutesPart: data.Sector2TimeMinutesPart,
		Sector3TimeMSPart:      data.Sector3TimeMSPart,
		Sector3TimeMinutesPart: data.Sector3TimeMinutesPart,
		LapValidBitFlags:       data.LapValidBitFlags,
	}, nil
}

func MapToTyreStintHistory(data *parsers.TyreStintHistoryData) (telentity.TyreStintHistoryData, error) {
	if data == nil {
		return telentity.TyreStintHistoryData{}, errors.New("mappers: nil TyreStintHistoryData")
	}

	return telentity.TyreStintHistoryData{
		EndLap:             data.EndLap,
		TyreActualCompound: consts.ActualTyreCompound[int16(data.TyreActualCompound)],
		TyreVisualCompound: consts.VisualTyreCompound[int16(data.TyreVisualCompound)],
	}, nil
}

func MapToSessionHistoryPacket(data *parsers.PacketSessionHistoryData) (telentity.SessionHistoryPacket, error) {
	if data == nil {
		return telentity.SessionHistoryPacket{}, errors.New("mappers: nil PacketSessionHistoryData")
	}

	header, err := MapToHeader(&data.Header)
	if err != nil {
		return telentity.SessionHistoryPacket{}, err
	}

	var laps [100]telentity.LapHistoryData
	for i := range data.LapHistoryData {
		lap, err := MapToLapHistory(&data.LapHistoryData[i])
		if err != nil {
			return telentity.SessionHistoryPacket{}, err
		}
		laps[i] = lap
	}

	var stints [8]telentity.TyreStintHistoryData
	for i := range data.TyreStintsHistoryData {
		stint, err := MapToTyreStintHistory(&data.TyreStintsHistoryData[i])
		if err != nil {
			return telentity.SessionHistoryPacket{}, err
		}
		stints[i] = stint
	}

	return telentity.SessionHistoryPacket{
		Header:                header,
		CarIdx:                data.CarIdx,
		NumLaps:               data.NumLaps,
		NumTyreStints:         data.NumTyreStints,
		BestLapTimeLapNum:     data.BestLapTimeLapNum,
		BestSector1LapNum:     data.BestSector1LapNum,
		BestSector2LapNum:     data.BestSector2LapNum,
		BestSector3LapNum:     data.BestSector3LapNum,
		LapHistoryData:        laps,
		TyreStintsHistoryData: stints,
	}, nil
}
