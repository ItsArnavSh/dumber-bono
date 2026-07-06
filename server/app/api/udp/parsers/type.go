package parsers

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Packet interface {
	PacketMotionData |
		PacketLapData |
		PacketSessionData |
		PacketEventData |
		PacketParticipantsData |
		PacketCarSetupData |
		PacketCarTelemetryData |
		PacketCarStatusData |
		PacketFinalClassificationData |
		PacketLobbyInfoData |
		PacketCarDamageData |
		PacketSessionHistoryData |
		PacketTyreSetsData |
		PacketLapPositionsData
}

func ParsePacket[P Packet](payload []byte) (*P, error) {
	var packet P
	reader := bytes.NewReader(payload)
	if err := binary.Read(reader, binary.LittleEndian, &packet); err != nil {
		return nil, fmt.Errorf("decode session packet: %w", err)
	}
	return &packet, nil
}
