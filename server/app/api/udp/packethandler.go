package udp

import (
	"dubmer-bono/app/api/udp/mappers"
	"dubmer-bono/app/api/udp/parsers"
	"fmt"
	"net"
)

func (u *UDPServer) handle_packet(addr net.Addr, payload []byte) {
	header, err := parsers.ParseHeader(payload)
	if err != nil {
		u.logger.Errorf("Error Parsing Header %v", err)
	}
	fmt.Println(header)

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
		_ = motion
		//Call Save From Here

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
		_ = session
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
		_ = lapData
		//Call Save From Here

	case 3:
		data, err := parsers.ParsePacket[parsers.PacketEventData](payload)
		if err != nil {
			u.logger.Errorf("Error parsing Event packet: %v", err)
			return
		}
		event, err := mappers.MapToEventData(data)
		if err != nil {
			u.logger.Errorf("Error mapping Event packet: %v", err)
			return
		}
		_ = event
		//Call Save From Here

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
		_ = participants
		//Call Save From Here

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
		_ = setup
		//Call Save From Here

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
		_ = telemetry
		//Call Save From Here

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
		_ = status
		//Call Save From Here

	case 8:
		// TODO: PacketFinalClassificationData has no mapper yet.
		_, err := parsers.ParsePacket[parsers.PacketFinalClassificationData](payload)
		if err != nil {
			u.logger.Errorf("Error parsing FinalClassification packet: %v", err)
			return
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
		_ = lobby
		//Call Save From Here

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
		_ = damage
		//Call Save From Here

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
		_ = history
		//Call Save From Here

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
		_ = tyreSets
		//Call Save From Here

	case 13:
		fmt.Println("Motion Ex") //Ignore (No Rig stuff needed here)

	case 14:
		fmt.Println("Time Trial") //Ignore

	case 15:
		fmt.Println("Lap Positions")
		// TODO: PacketLapPositionsData has no mapper yet.
		_, err := parsers.ParsePacket[parsers.PacketLapPositionsData](payload)
		if err != nil {
			u.logger.Errorf("Error parsing LapPositions packet: %v", err)
			return
		}

	default:
		fmt.Printf("Unknown PacketID: %d\n", header.PacketID)
	}
}
