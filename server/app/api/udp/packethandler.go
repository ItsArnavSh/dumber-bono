package udp

import (
	"dubmer-bono/app/api/udp/mappers"
	"dubmer-bono/app/api/udp/parsers"
	"fmt"
	"net"
)

// VerboseLogging toggles full struct dumps for every parsed/mapped packet.
// Set to true to log the packet type and all its contents; false to stay quiet
// except for errors.
var VerboseLogging = false

func (u *UDPServer) handle_packet(addr net.Addr, payload []byte) {
	header, err := parsers.ParseHeader(payload)

	if err != nil {
		u.logger.Errorf("Error Parsing Header %v", err)
	}
	if u.throttle.Allow() {
		u.service.IngestHeader(header)
	}

	switch header.PacketID {
	case 0:
		data, err := parsers.ParsePacket[parsers.PacketMotionData](payload)
		if err != nil {
			u.logger.Errorf("Error parsing Motion packet: %v", err)
			return
		}
		motion, err := mappers.MaptoMotionPacket(data)
		if err != nil {
			u.logger.Errorf("Error mapping Motion packet: %v", err)
			return
		}
		if VerboseLogging {
			u.logger.Infof("Motion packet: %+v", motion)
		}
		u.service.IngestMotionPacket(motion)

	case 1:
		data, err := parsers.ParsePacket[parsers.PacketSessionData](payload)
		if err != nil {
			u.logger.Errorf("Error parsing Session packet: %v", err)
			return
		}
		session, err := mappers.MapToSessionData(data)
		if err != nil {
			u.logger.Errorf("Error mapping Session packet: %v", err)
			return
		}
		if VerboseLogging {
			u.logger.Infof("Session packet: %+v", session)
		}
		u.service.IngestSessionPacket(session)
		//Call Save From Here

	case 2:
		data, err := parsers.ParsePacket[parsers.PacketLapData](payload)
		if err != nil {
			u.logger.Errorf("Error parsing LapData packet: %v", err)
			return
		}
		lapData, err := mappers.MapToLapDataPacket(data)
		if err != nil {
			u.logger.Errorf("Error mapping LapData packet: %v", err)
			return
		}
		if VerboseLogging {
			u.logger.Infof("LapData packet: %+v", lapData)
		}
		u.service.IngestLapPacket(lapData)

	case 3:
		data, err := parsers.ParseEventPacket(payload)
		if err != nil {
			u.logger.Errorf("Error parsing Event packet: %v", err)
			return
		}
		event, err := mappers.MapToEventData(data)
		if err != nil {
			u.logger.Errorf("Error mapping Event packet: %v", err)
			return
		}
		if VerboseLogging {
			u.logger.Infof("Event packet: %+v", event)
		}

		u.service.IngestEventPacket(event)

	case 4:
		data, err := parsers.ParsePacket[parsers.PacketParticipantsData](payload)
		if err != nil {
			u.logger.Errorf("Error parsing Participants packet: %v", err)
			return
		}
		participants, err := mappers.MapToParticipantsPacket(data)
		if err != nil {
			u.logger.Errorf("Error mapping Participants packet: %v", err)
			return
		}
		if VerboseLogging {
			u.logger.Infof("Participants packet: %+v", participants)
		}

		u.service.IngestParticipantPacket(participants)

	case 5:
		data, err := parsers.ParsePacket[parsers.PacketCarSetupData](payload)
		if err != nil {
			u.logger.Errorf("Error parsing CarSetup packet: %v", err)
			return
		}
		setup, err := mappers.MapToCarSetupPacket(data)
		if err != nil {
			u.logger.Errorf("Error mapping CarSetup packet: %v", err)
			return
		}
		if VerboseLogging {
			u.logger.Infof("CarSetup packet: %+v", setup)
		}

		u.service.IngestCarSetupPacket(setup)

	case 6:
		data, err := parsers.ParsePacket[parsers.PacketCarTelemetryData](payload)
		if err != nil {
			u.logger.Errorf("Error parsing CarTelemetry packet: %v", err)
			return
		}
		telemetry, err := mappers.MapToCarTelemetryPacket(data)
		if err != nil {
			u.logger.Errorf("Error mapping CarTelemetry packet: %v", err)
			return
		}
		if VerboseLogging {
			u.logger.Infof("CarTelemetry packet: %+v", telemetry)
		}

		u.service.IngestTelemetryPacket(telemetry)

	case 7:
		data, err := parsers.ParsePacket[parsers.PacketCarStatusData](payload)
		if err != nil {
			u.logger.Errorf("Error parsing CarStatus packet: %v", err)
			return
		}
		status, err := mappers.MapToCarStatusPacket(data)
		if err != nil {
			u.logger.Errorf("Error mapping CarStatus packet: %v", err)
			return
		}
		if VerboseLogging {
			u.logger.Infof("CarStatus packet: %+v", status)
		}

		u.service.IngestCarStatusPacket(status)

	case 8:
		// TODO: PacketFinalClassificationData has no mapper yet.
		data, err := parsers.ParsePacket[parsers.PacketFinalClassificationData](payload)
		if err != nil {
			u.logger.Errorf("Error parsing FinalClassification packet: %v", err)
			return
		}
		if VerboseLogging {
			u.logger.Infof("FinalClassification packet (raw, no mapper yet): %+v", data)
		}

	case 9:
		data, err := parsers.ParsePacket[parsers.PacketLobbyInfoData](payload)
		if err != nil {
			u.logger.Errorf("Error parsing LobbyInfo packet: %v", err)
			return
		}
		lobby, err := mappers.MapToLobbyInfoPacket(data)
		if err != nil {
			u.logger.Errorf("Error mapping LobbyInfo packet: %v", err)
			return
		}
		if VerboseLogging {
			u.logger.Infof("LobbyInfo packet: %+v", lobby)
		}

		u.service.IngestLobbyInfoPacket(lobby)

	case 10:
		data, err := parsers.ParsePacket[parsers.PacketCarDamageData](payload)
		if err != nil {
			u.logger.Errorf("Error parsing CarDamage packet: %v", err)
			return
		}
		damage, err := mappers.MapToCarDamagePacket(data)
		if err != nil {
			u.logger.Errorf("Error mapping CarDamage packet: %v", err)
			return
		}
		if VerboseLogging {
			u.logger.Infof("CarDamage packet: %+v", damage)
		}

		u.service.IngestCarDamagePacket(damage)

	case 11:
		data, err := parsers.ParsePacket[parsers.PacketSessionHistoryData](payload)
		if err != nil {
			u.logger.Errorf("Error parsing SessionHistory packet: %v", err)
			return
		}
		history, err := mappers.MapToSessionHistoryPacket(data)
		if err != nil {
			u.logger.Errorf("Error mapping SessionHistory packet: %v", err)
			return
		}
		if VerboseLogging {
			u.logger.Infof("SessionHistory packet: %+v", history)
		}

		u.service.IngestSessionHistoryPacket(history)

	case 12:
		data, err := parsers.ParsePacket[parsers.PacketTyreSetsData](payload)
		if err != nil {
			u.logger.Errorf("Error parsing TyreSets packet: %v", err)
			return
		}
		tyreSets, err := mappers.MapToTyreSetPacket(data)
		if err != nil {
			u.logger.Errorf("Error mapping TyreSets packet: %v", err)
			return
		}
		if VerboseLogging {
			u.logger.Infof("TyreSets packet: %+v", tyreSets)
		}

		u.service.IngestTyreSetPacket(tyreSets)

	case 13:
		if VerboseLogging {
			u.logger.Infof("Motion Ex packet received (ignored, no rig support needed)")
		}

	case 14:
		if VerboseLogging {
			u.logger.Infof("Time Trial packet received (ignored)")
		}

	case 15:
		// TODO: PacketLapPositionsData has no mapper yet.
		data, err := parsers.ParsePacket[parsers.PacketLapPositionsData](payload)
		if err != nil {
			u.logger.Errorf("Error parsing LapPositions packet: %v", err)
			return
		}
		if VerboseLogging {
			u.logger.Infof("LapPositions packet (raw, no mapper yet): %+v", data)
		}

	default:
		fmt.Printf("Unknown PacketID: %d\n", header.PacketID)
	}
}
