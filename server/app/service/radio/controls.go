package radio

import "dubmer-bono/app/types/entity"

func (r *Service) HandleHotKeyEvents() {
	for hkey := range r.hotkey_chan {
		switch hkey {
		case entity.RADIO_PRESS:
		case entity.RADIO_RELEASE:
		case entity.COPY_AFFIRMATION:
		case entity.MUTE_TOGGLE:
			r.muted = !r.muted
			if r.muted {
				r.logger.Info("Radio Muted")
			} else {
				r.logger.Info("Radio Unmuted")
			}
		}
	}
}
