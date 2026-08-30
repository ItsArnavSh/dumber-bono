package api

import (
	"context"
	"dubmer-bono/app/api/hotkeys"
	"dubmer-bono/app/api/udp"
	"dubmer-bono/app/service"
	"dubmer-bono/app/types/entity"
	"dubmer-bono/app/utility/errors"

	"go.uber.org/zap"
)

type Server struct {
	logger   *zap.SugaredLogger
	services Services
}

func NewServer(ctx context.Context, logger *zap.SugaredLogger, path string) (Server, error) {
	bg_errs := errors.NewErrorHandler(logger)

	//init the Repo(cache, db, analytics server)
	repo, err := service.NewRepository(ctx, path)
	if err != nil {
		logger.Errorf("error initializing server: %w", err)
		return Server{}, err
	}

	//init the Hotkey Listner
	hkey_listner, err := hotkeys.NewHotKeyListner()
	if err != nil {
		logger.Errorf("error setting up hotkeys listner", err)
		return Server{}, err
	}
	eventListner := make(chan entity.HotKeyEvent)
	if err := hkey_listner.InitHandler(eventListner); err != nil {
		logger.Errorf("error initializing hotkeys listener: %v", err)
		return Server{}, err
	}
	go func() {
		if err := hkey_listner.StartListening(ctx); err != nil && err != context.Canceled {
			bg_errs.NewError("HotKey listener crashed", err, errors.FATAL)
		}
	}()

	//Init the Services
	services, err := InitServices(ctx, path, logger, repo, eventListner)
	if err != nil {
		logger.Errorf("error initializing services: %w", err)
		return Server{}, err
	}

	//Init the UDP Server
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

func (s *Server) CloseServer(ctx context.Context) {

}
