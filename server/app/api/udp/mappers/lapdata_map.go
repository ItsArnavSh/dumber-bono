package mappers

import (
	"errors"

	"dubmer-bono/app/api/udp/parsers"
	"dubmer-bono/app/types/entity"
	"dubmer-bono/app/types/entity/consts"
	telentity "dubmer-bono/app/types/entity/tel-entity"
)

func MapToLapData(data *parsers.LapData) (telentity.LapData, error) {
	if data == nil {
		return telentity.LapData{}, errors.New("mappers: nil LapData")
	}

	return telentity.LapData{
		LastLapTime:    data.LastLapTimeInMS,
		CurrentLapTime: data.CurrentLapTimeInMS,
		SectorTimes: entity.SectorWise{
			Sector1: entity.LapTimeStamp{Minute: data.Sector1TimeMinutesPart, MS: data.Sector1TimeMSPart},
			Sector2: entity.LapTimeStamp{Minute: data.Sector2TimeMinutesPart, MS: data.Sector2TimeMSPart},
			// Sector3: no source field in parsers.LapData; left zero-valued.
		},
		DeltaToFront: entity.LapTimeStamp{
			Minute: data.DeltaToCarInFrontMinutesPart,
			MS:     data.DeltaToCarInFrontMSPart,
		},
		DeltaToRaceLeader: entity.LapTimeStamp{
			Minute: data.DeltaToRaceLeaderMinutesPart,
			MS:     data.DeltaToRaceLeaderMSPart,
		},
		LapDistance:                 data.LapDistance,
		TotalDistanceInSession:      data.TotalDistance,
		SafetyCarDelta:              data.SafetyCarDelta,
		CarPosition:                 data.CarPosition,
		CurrentLapNum:               data.CurrentLapNum,
		PitStatus:                   consts.PitStatus[int16(data.PitStatus)],
		NumPitStops:                 data.NumPitStops,
		Sector:                      consts.Sector[int16(data.Sector)],
		CurrentLapInvalid:           data.CurrentLapInvalid != 0,
		Penalties:                   data.Penalties,
		TotalWarnings:               data.TotalWarnings,
		CornerCuttingWarnings:       data.CornerCuttingWarnings,
		NumUnservedDriveThroughPens: data.NumUnservedDriveThroughPens,
		NumUnservedStopGoPens:       data.NumUnservedStopGoPens,
		GridPosition:                data.GridPosition,
		DriverStatus:                consts.DriverStatus[int16(data.DriverStatus)],
		ResultStatus:                consts.ResultStatus[int16(data.ResultStatus)],
		PitLaneTimerActive:          data.PitLaneTimerActive != 0,
		PitLaneTimeInLaneInMS:       data.PitLaneTimeInLaneInMS,
		PitStopTimerInMS:            data.PitStopTimerInMS,
		PitStopShouldServePen:       data.PitStopShouldServePen != 0,
		SpeedTrapFastestSpeed:       data.SpeedTrapFastestSpeed,
		SpeedTrapFastestLap:         data.SpeedTrapFastestLap,
	}, nil
}

func MapToLapDataPacket(data *parsers.PacketLapData) (telentity.LapDataPacket, error) {
	if data == nil {
		return telentity.LapDataPacket{}, errors.New("mappers: nil PacketLapData")
	}

	header, err := MapToHeader(&data.Header)
	if err != nil {
		return telentity.LapDataPacket{}, err
	}

	var laps [22]telentity.LapData
	for i := range data.LapData {
		lap, err := MapToLapData(&data.LapData[i])
		if err != nil {
			return telentity.LapDataPacket{}, err
		}
		laps[i] = lap
	}

	return telentity.LapDataPacket{
		Header:  header,
		LapData: laps,
	}, nil
}
