package radio

import (
	"context"
	"dubmer-bono/app/types/entity"
)

func (r *Service) HandleHotKeyEvents(ctx context.Context) {
	for hkey := range r.hotkey_chan {
		switch hkey {
		case entity.RADIO_PRESS:
			r.stt.StartMessageRec(ctx)
			r.logger.Infof("Started Listening")
		case entity.RADIO_RELEASE:
			message, err := r.stt.EndMessageRec(ctx)
			if err != nil {
				r.logger.Error(err)
			} else {
				r.logger.Infof("Message: %s", message)
			}

		case entity.COPY_AFFIRMATION:
		//Future Use

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
