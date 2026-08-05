package types

import (
	"context"
	"dubmer-bono/app/types/entity"
)

type HotKeyHandler interface {
	InitHandler(eventChan chan<- entity.HotKeyEvent) error
	StartListening(ctx context.Context) error
	StopListening(ctx context.Context) error
}
