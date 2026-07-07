package mappers

import (
	"errors"

	"dubmer-bono/app/api/udp/parsers"
	"dubmer-bono/app/types/entity/consts"
	telentity "dubmer-bono/app/types/entity/tel-entity"
)

func zoneFlagString(f int8) string {
	switch f {
	case -1:
		return "invalid/unknown"
	case 0:
		return "none"
	case 1:
		return "green"
	case 2:
		return "blue"
	case 3:
		return "yellow"
	default:
		return "unknown"
	}
}

func MapToMarshalZone(data *parsers.MarshalZone) (telentity.MarshalZone, error) {
	if data == nil {
		return telentity.MarshalZone{}, errors.New("mappers: nil MarshalZone")
	}
	return telentity.MarshalZone{
		ZoneStart: data.ZoneStart,
		ZoneFlag:  zoneFlagString(data.ZoneFlag),
	}, nil
}

func MapToWeatherForecastSample(data *parsers.WeatherForecastSample) (telentity.WeatherForecastSample, error) {
	if data == nil {
		return telentity.WeatherForecastSample{}, errors.New("mappers: nil WeatherForecastSample")
	}
	return telentity.WeatherForecastSample{
		SessionType:            consts.SessionTypes[uint16(data.SessionType)],
		TimeOffset:             data.TimeOffset,
		Weather:                consts.WeatherForecastTypes[uint16(data.Weather)],
		TrackTemperature:       data.TrackTemperature,
		TrackTemperatureChange: consts.TempChange[int16(data.TrackTemperatureChange)],
		AirTemperature:         data.AirTemperature,
		AirTemperatureChange:   consts.TempChange[int16(data.AirTemperatureChange)],
		RainPercentage:         data.RainPercentage,
	}, nil
}

func MapToSessionData(data *parsers.PacketSessionData) (telentity.PacketSessionData, error) {
	if data == nil {
		return telentity.PacketSessionData{}, errors.New("mappers: nil PacketSessionData")
	}

	header, err := MapToHeader(&data.Header)
	if err != nil {
		return telentity.PacketSessionData{}, err
	}

	var marshalZones [21]telentity.MarshalZone
	for i := range data.MarshalZones {
		mz, err := MapToMarshalZone(&data.MarshalZones[i])
		if err != nil {
			return telentity.PacketSessionData{}, err
		}
		marshalZones[i] = mz
	}

	numSamples := min(int(data.NumWeatherForecastSamples), len(data.WeatherForecastSamples))
	weatherSamples := make([]telentity.WeatherForecastSample, 0, numSamples)
	for i := range numSamples {
		ws, err := MapToWeatherForecastSample(&data.WeatherForecastSamples[i])
		if err != nil {
			return telentity.PacketSessionData{}, err
		}
		weatherSamples = append(weatherSamples, ws)
	}

	trackName := "unknown"
	if data.TrackId >= 0 {
		trackName = consts.TrackIDs[uint16(data.TrackId)]
	}

	return telentity.PacketSessionData{
		Header:                  header,
		Weather:                 consts.WeatherForecastTypes[uint16(data.Weather)],
		TrackTemperature:        data.TrackTemperature,
		AirTempterature:         data.AirTemperature,
		TotalLaps:               int8(data.TotalLaps),
		SessionType:             consts.SessionTypes[uint16(data.SessionType)],
		Track:                   trackName,
		Formula:                 consts.Formula[int16(data.Formula)],
		SessionTimeLeft:         data.SessionTimeLeft,
		SessionDuration:         data.SessionDuration,
		PitSpeedLimit:           data.PitSpeedLimit,
		GamePaused:              data.GamePaused != 0,
		NumMarshalZones:         data.NumMarshalZones,
		MarshalZones:            marshalZones,
		SafetyCarStatus:         consts.SafetyCarStatus[int16(data.SafetyCarStatus)],
		NetworkGame:             consts.NetworkGame[int16(data.NetworkGame)],
		WeatherForecastSamples:  weatherSamples,
		PitStopWindowIdealLap:   data.PitStopWindowIdealLap,
		PitStopWindowLatestLap:  data.PitStopWindowLatestLap,
		PitStopRejoinPosition:   data.PitStopRejoinPosition,
		GameMode:                data.GameMode,
		RuleSet:                 data.RuleSet,
		TimeOfDay:               data.TimeOfDay,
		SessionLength:           consts.SessionLength[int16(data.SessionLength)],
		Sector2LapDistanceStart: data.Sector2LapDistanceStart,
		Sector3LapDistanceStart: data.Sector3LapDistanceStart,
	}, nil
}
