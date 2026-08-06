package mappers

import (
	"dubmer-bono/app/api/udp/parsers"
	telentity "dubmer-bono/app/types/entity/tel-entity"
	"errors"
)

func MapToHeader(data *parsers.PacketHeader) (telentity.UDPHeader, error) {
	if data == nil {
		return telentity.UDPHeader{}, errors.New("mappers: nil PacketHeader")
	}

	return telentity.UDPHeader{
		PacketFormat:            data.PacketFormat,
		GameYear:                data.GameYear,
		GameMajorVersion:        data.GameMajorVersion,
		GameMinorVersion:        data.GameMinorVersion,
		PacketVersion:           data.PacketVersion,
		PacketID:                data.PacketID,
		SessionUID:              data.SessionUID,
		SessionTime:             data.SessionTime,
		FrameIdentifier:         data.FrameIdentifier,
		PlayerCarIndex:          data.PlayerCarIndex,
		SecondaryPlayerCarIndex: data.SecondaryPlayerCarIndex,
	}, nil
}
