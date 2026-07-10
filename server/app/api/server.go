package api

import (
	"context"
	"dubmer-bono/app/api/udp"

	"go.uber.org/zap"
)

type Server struct {
	logger   *zap.SugaredLogger
	services Services
}

func NewServer(ctx context.Context, logger *zap.SugaredLogger, path string) (Server, error) {
	services, err := InitServices(ctx, path, logger)
	if err != nil {
		logger.Errorf("error initializing services: %w", err)
		return Server{}, err
	}

	err = udp.ListenUDP(ctx, logger, 4345, services.ingestion)
	if err != nil {
		logger.Errorf("error setting up UDP server: %w", err)
		return Server{}, err
	}

	return Server{
		logger:   logger,
		services: services,
	}, nil
}

func (s *Server) StartServer(ctx context.Context) {
	select {}
}
