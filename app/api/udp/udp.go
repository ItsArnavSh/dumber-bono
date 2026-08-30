package udp

import (
	"context"
	"dubmer-bono/app/types"
	"dubmer-bono/app/types/entity/consts"
	"dubmer-bono/app/utility"
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
)

type UDPServer struct {
	throttle *utility.Throttler
	port     uint16
	logger   *zap.SugaredLogger
	conn     net.PacketConn
	service  types.Ingestion
}

func (u *UDPServer) listenUDP(ctx context.Context) error {
	defer func() {
		if err := u.conn.Close(); err != nil {
			u.logger.Errorf("Error closing UDP connection: %v", err)
		}
	}()
	for {
		buf := make([]byte, consts.UDP_PACKET_SIZE)
		n, addr, err := u.conn.ReadFrom(buf)
		if err != nil {
			//Perhaps should not crash here, we'll see
			u.logger.Errorf("Error reading buffer: %w", err)
			continue
		}
		go u.handle_packet(addr, buf[:n])
	}
}

func ListenUDP(ctx context.Context, logger *zap.SugaredLogger, port uint16, service types.Ingestion) error {
	server := UDPServer{
		port:     port,
		logger:   logger,
		service:  service,
		throttle: &utility.Throttler{Interval: time.Second},
	}
	pc, err := (&net.ListenConfig{}).ListenPacket(ctx, "udp", fmt.Sprintf(":%d", server.port))
	if err != nil {
		server.logger.Errorf("UDP server could not be initialised: %w", err)
		return err
	}

	server.conn = pc
	go func() {
		if err := server.listenUDP(ctx); err != nil {
			server.logger.Errorf("UDP listener stopped: %v", err)
		}
	}()

	server.logger.Info("Started The UDP Server")
	return nil
}
