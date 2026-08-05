package general

import (
	"context"
	"dubmer-bono/app/types/entity"

	"fmt"
	"log"

	"golang.design/x/hotkey"
)

type CrossPlatformHotkeys struct {
	events  chan<- entity.HotKeyEvent
	radioHK *hotkey.Hotkey
	copyHK  *hotkey.Hotkey
	muteHK  *hotkey.Hotkey
}

func (h *CrossPlatformHotkeys) InitHandler(eventChan chan<- entity.HotKeyEvent) error {
	h.events = eventChan

	// Single key hotkeys (empty modifier list means just the raw key)
	h.radioHK = hotkey.New([]hotkey.Modifier{}, hotkey.KeyR)
	h.copyHK = hotkey.New([]hotkey.Modifier{}, hotkey.KeyC)
	h.muteHK = hotkey.New([]hotkey.Modifier{}, hotkey.KeyM)

	return nil
}

func (h *CrossPlatformHotkeys) StartListening(ctx context.Context) error {
	if err := h.radioHK.Register(); err != nil {
		return fmt.Errorf("failed to register R key: %w", err)
	}
	if err := h.copyHK.Register(); err != nil {
		h.radioHK.Unregister()
		return fmt.Errorf("failed to register C key: %w", err)
	}
	if err := h.muteHK.Register(); err != nil {
		h.radioHK.Unregister()
		h.copyHK.Unregister()
		return fmt.Errorf("failed to register M key: %w", err)
	}

	log.Println("Single-key global hotkeys registered. Listening...")

	go h.listenHotkey(ctx, h.radioHK, entity.RADIO_PRESS, entity.RADIO_RELEASE)
	go h.listenHotkeyPressOnly(ctx, h.copyHK, entity.COPY_AFFIRMATION)
	go h.listenHotkeyPressOnly(ctx, h.muteHK, entity.MUTE_TOGGLE)

	<-ctx.Done()
	return nil
}

func (h *CrossPlatformHotkeys) StopListening(ctx context.Context) error {
	var errs []error

	if h.radioHK != nil {
		if err := h.radioHK.Unregister(); err != nil {
			errs = append(errs, err)
		}
	}
	if h.copyHK != nil {
		if err := h.copyHK.Unregister(); err != nil {
			errs = append(errs, err)
		}
	}
	if h.muteHK != nil {
		if err := h.muteHK.Unregister(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors unregistering hotkeys: %v", errs)
	}
	return nil
}

func (h *CrossPlatformHotkeys) listenHotkey(ctx context.Context, hk *hotkey.Hotkey, pressEvent, releaseEvent entity.HotKeyEvent) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-hk.Keydown():
			select {
			case h.events <- pressEvent:
			default:
			}
		case <-hk.Keyup():
			select {
			case h.events <- releaseEvent:
			default:
			}
		}
	}
}

func (h *CrossPlatformHotkeys) listenHotkeyPressOnly(ctx context.Context, hk *hotkey.Hotkey, event entity.HotKeyEvent) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-hk.Keydown():
			select {
			case h.events <- event:
			default:
			}
		}
	}
}
