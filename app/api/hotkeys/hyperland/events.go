package hyperland

import (
	"dubmer-bono/app/types/entity"
	"log"
	"syscall"
	"unsafe"
)

func (h *HyperlandHotkeys) monitorKeys() error {
	eventSize := unsafe.Sizeof(InputEvent{})
	buf := make([]byte, eventSize)

	for {
		n, err := h.deviceFile.Read(buf)
		if err != nil {
			log.Printf("Error reading from %s: %v\n", h.devicePath, err)
			break
		}

		if n < int(eventSize) {
			continue
		}

		event := (*InputEvent)(unsafe.Pointer(&buf[0]))

		// Type 1 = keyboard event (EV_KEY)
		if event.Type == 1 {
			// event.Value: 0=release, 1=press, 2=repeat
			// Skip repeat events (value=2)
			if event.Value == 2 {
				continue
			}

			var hotKeyEvent entity.HotKeyEvent
			matched := true

			switch event.Code {
			case 19: // R key
				if event.Value == 1 {
					hotKeyEvent = entity.RADIO_PRESS
				} else {
					hotKeyEvent = entity.RADIO_RELEASE
				}
			case 46: // C key (Only trigger on press to match typical shortcuts)
				if event.Value == 1 {
					hotKeyEvent = entity.COPY_AFFIRMATION
				} else {
					continue // Skip release for C
				}
			case 50: // M key (Only trigger on press)
				if event.Value == 1 {
					hotKeyEvent = entity.MUTE_TOGGLE
				} else {
					continue // Skip release for M
				}
			default:
				matched = false
			}

			if matched {
				h.events <- hotKeyEvent
			}
		}
	}
	return nil
}

// InputEvent represents a single input event from the kernel
type InputEvent struct {
	Time  syscall.Timeval
	Type  uint16
	Code  uint16
	Value int32
}
