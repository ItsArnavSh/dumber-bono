package mappers

import (
	"errors"

	"dubmer-bono/app/api/udp/parsers"
	telentity "dubmer-bono/app/types/entity/tel-entity"
)

func yourTelemetryString(v uint8) string {
	if v == 1 {
		return "public"
	}
	return "restricted"
}

func MapToParticipant(data *parsers.ParticipantData) (telentity.ParticipantData, error) {
	if data == nil {
		return telentity.ParticipantData{}, errors.New("mappers: nil ParticipantData")
	}

	return telentity.ParticipantData{
		DriverId:      data.DriverId,
		TeamId:        data.TeamId,
		MyTeam:        data.MyTeam != 0,
		RaceNumber:    data.RaceNumber,
		Nationality:   data.Nationality,
		Name:          bytesToName(data.Name),
		YourTelemetry: yourTelemetryString(data.YourTelemetry),
	}, nil
}

func MapToParticipantsPacket(data *parsers.PacketParticipantsData) (telentity.PacketParticipantsData, error) {
	if data == nil {
		return telentity.PacketParticipantsData{}, errors.New("mappers: nil PacketParticipantsData")
	}

	header, err := MapToHeader(&data.Header)
	if err != nil {
		return telentity.PacketParticipantsData{}, err
	}

	var participants [22]telentity.ParticipantData
	for i := range data.Participants {
		p, err := MapToParticipant(&data.Participants[i])
		if err != nil {
			return telentity.PacketParticipantsData{}, err
		}
		participants[i] = p
	}

	return telentity.PacketParticipantsData{
		Header:        header,
		NumActiveCars: data.NumActiveCars,
		Participants:  participants,
	}, nil
}
