package errors

import (
	"log"

	"go.uber.org/zap"
)

type ErrorEvent struct {
	Err         error
	Consequence ErrorCons
}

type ErrorCons int

const (
	FATAL ErrorCons = iota
	LOG             //Just logs does not crash the main thread
)

type GoErrorHandler struct {
	logger *zap.SugaredLogger
}

func NewErrorHandler(logger *zap.SugaredLogger) GoErrorHandler {
	return GoErrorHandler{logger: logger}
}

func (g *GoErrorHandler) NewError(message string, err error, consequence ErrorCons) {
	switch consequence {
	case LOG:
		g.logger.Errorf("%s: %v", message, err)
	case FATAL:
		log.Fatalf("%s:%v", message, err)
	}
}
