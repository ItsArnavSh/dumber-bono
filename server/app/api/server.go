package api

import (
	"dubmer-bono/app/api/udp"

	"go.uber.org/zap"
)

type Server struct {
	logger *zap.SugaredLogger
	udp    udp.UDPServer
}
