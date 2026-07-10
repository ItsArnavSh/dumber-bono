package types

import (
	"dubmer-bono/app/api/udp/parsers"
	telentity "dubmer-bono/app/types/entity/tel-entity"
)

type Ingestion interface {
	IngestHeader(payload *parsers.PacketHeader) error
	IngestMotionPacket(payload telentity.MotionPacket)
	IngestSessionPacket(payload telentity.PacketSessionData)
	IngestLapPacket(payload telentity.LapDataPacket)
	IngestEventPacket(payload telentity.PacketEventData)
	IngestParticipantPacket(payload telentity.PacketParticipantsData)
	IngestCarSetupPacket(payload telentity.PacketCarSetupData)
	IngestTelemetryPacket(payload telentity.PacketCarTelemetryData)
	IngestCarStatusPacket(payload telentity.PacketCarStatusData)
	IngestLobbyInfoPacket(payload telentity.PacketLobbyInfoData)
	IngestCarDamagePacket(payload telentity.PacketCarDamageData)
	IngestSessionHistoryPacket(payload telentity.SessionHistoryPacket)
	IngestTyreSetPacket(payload telentity.TyreSetPacket)
}

type Monitor interface {
	GetPressure() int
}

type Radio interface {
}

type Agent interface {
}
