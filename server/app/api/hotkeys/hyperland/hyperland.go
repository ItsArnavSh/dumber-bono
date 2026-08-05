package hyperland

import (
	"context"
	"dubmer-bono/app/types/entity"
	"fmt"
	"log"
	"os"
)

// This is for testing purposes only in my hyperland setup, the windows version will be used in actual app
//

type HyperlandHotkeys struct {
	events     chan<- entity.HotKeyEvent
	devicePath string
	deviceFile *os.File // To listen to keyboard events
}

func (h *HyperlandHotkeys) InitHandler(eventChan chan<- entity.HotKeyEvent) error {
	h.events = eventChan
	h.devicePath = os.Getenv("KEYBOARD_DEVICE")
	if h.devicePath == "" {
		return fmt.Errorf("KEYBOARD_DEVICE env not set")
	}
	return nil
}

func (h *HyperlandHotkeys) StartListening(ctx context.Context) error {
	file, err := os.Open(h.devicePath)
	if err != nil {
		log.Printf("Error opening %s: %v\n", h.devicePath, err)
		return err
	}
	h.deviceFile = file
	return h.monitorKeys()
}

func (h *HyperlandHotkeys) StopListening(ctx context.Context) error {
	return h.deviceFile.Close()
}
