package udp

import (
	"context"
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
	return nil
}
func (u *UDPServer) handle_packet(addr net.Addr, payload []byte) {}
