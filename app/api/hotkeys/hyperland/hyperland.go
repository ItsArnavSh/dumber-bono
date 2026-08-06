package hyperland

import (
	"context"
	"dubmer-bono/app/types/entity"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	h.devicePath = findKeyboard()
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
func findKeyboard() string {
	byIDDir := "/dev/input/by-id/"
	entries, err := os.ReadDir(byIDDir)
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		name := entry.Name()
		// Look for keyboard-related entries that do not contain "Mouse"
		if strings.Contains(name, "event-kbd") && !strings.Contains(name, "Mouse") {
			symlinkPath := filepath.Join(byIDDir, name)

			// Resolve the symlink (e.g., ../event5)
			realPath, err := os.Readlink(symlinkPath)
			if err != nil {
				continue
			}

			// Clean and resolve it relative to the directory to get the full absolute path
			absolutePath := filepath.Clean(filepath.Join(byIDDir, realPath))
			return absolutePath
		}
	}

	return ""
}
