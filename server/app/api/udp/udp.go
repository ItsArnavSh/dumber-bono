package udp

import (
	"context"
	"dubmer-bono/app/api/udp/mappers"
	"dubmer-bono/app/types/entity/consts"
	"fmt"
	"net"

	"go.uber.org/zap"
)

type UDPServer struct {
	port   uint16
	logger *zap.SugaredLogger
	conn   net.PacketConn
}

func (u *UDPServer) listenUDP(ctx context.Context) error {
	defer u.conn.Close()
	for {
		buf := make([]byte, consts.UDP_PACKET_SIZE)
		n, addr, err := u.conn.ReadFrom(buf)
		if err != nil {
			//Perhaps should not crash here, we'll see
			u.logger.Errorf("Error reading buffer: %w", err)
		}
		go u.handle_packet(addr, buf[:n])
	}
}

func NewUDPServer(ctx context.Context, logger *zap.SugaredLogger, port uint16) error {
	server := UDPServer{
		port:   port,
		logger: logger,
	}
	pc, err := net.ListenPacket("udp", fmt.Sprintf(":%d", server.port))
	if err != nil {
		server.logger.Errorf("UDP server could not be initialised: %w", err)
		return err
	}

	server.conn = pc
	go server.listenUDP(ctx)

	server.logger.Info("Started The UDP Server")
	return nil
}
func (u *UDPServer) handle_packet(addr net.Addr, payload []byte) {
	header, err := mappers.ParseHeader(payload)
	if err != nil {
		u.logger.Errorf("Error Parsing Header %v", err)
	}
	fmt.Println(header)
	switch header.PacketID {
	case 0:
		fmt.Println("Motion")
	case 1:
		fmt.Println("Session")
	case 2:
		fmt.Println("Lap Data")
	case 3:
		fmt.Println("Event")
	case 4:
		fmt.Println("Participants")
	case 5:
		fmt.Println("Car Setups")
	case 6:
		fmt.Println("Car Telemetry")
	case 7:
		fmt.Println("Car Status")
	case 8:
		fmt.Println("Final Classification")
	case 9:
		fmt.Println("Lobby Info")
	case 10:
		fmt.Println("Car Damage")
	case 11:
		fmt.Println("Session History")
	case 12:
		fmt.Println("Tyre Sets")
	case 13:
		fmt.Println("Motion Ex")
	case 14:
		fmt.Println("Time Trial")
	case 15:
		fmt.Println("Lap Positions")
	default:
		fmt.Printf("Unknown PacketID: %d\n", header.PacketID)
	}
}
