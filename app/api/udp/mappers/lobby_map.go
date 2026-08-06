package mappers

import (
	"bytes"
	"errors"

	"dubmer-bono/app/api/udp/parsers"
	"dubmer-bono/app/types/entity/consts"
	telentity "dubmer-bono/app/types/entity/tel-entity"
)

func bytesToName(name [32]byte) string {
	n := bytes.IndexByte(name[:], 0)
	if n == -1 {
		n = len(name)
	}
	return string(name[:n])
}

func MapToLobbyInfo(data *parsers.LobbyInfoData) (telentity.LobbyInfoData, error) {
	if data == nil {
		return telentity.LobbyInfoData{}, errors.New("mappers: nil LobbyInfoData")
	}

	return telentity.LobbyInfoData{
		AiControlled: data.AiControlled != 0,
		TeamId:       consts.TeamIDs[uint16(data.TeamId)],
		Nationality:  consts.NationalityIDs[uint16(data.Nationality)],
		Name:         bytesToName(data.Name),
		CarNumber:    data.CarNumber,
	}, nil
}

func MapToLobbyInfoPacket(data *parsers.PacketLobbyInfoData) (telentity.PacketLobbyInfoData, error) {
	if data == nil {
		return telentity.PacketLobbyInfoData{}, errors.New("mappers: nil PacketLobbyInfoData")
	}

	header, err := MapToHeader(&data.Header)
	if err != nil {
		return telentity.PacketLobbyInfoData{}, err
	}

	var players [22]telentity.LobbyInfoData
	for i := range data.LobbyPlayers {
		player, err := MapToLobbyInfo(&data.LobbyPlayers[i])
		if err != nil {
			return telentity.PacketLobbyInfoData{}, err
		}
		players[i] = player
	}

	return telentity.PacketLobbyInfoData{
		Header:       header,
		NumPlayers:   data.NumPlayers,
		LobbyPlayers: players,
	}, nil
}
