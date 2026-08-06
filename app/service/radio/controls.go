package radio

import (
	"context"
	"dubmer-bono/app/types/entity"
	"io"
	"time"

	"dubmer-bono/app/service/internal/llm"
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
				continue
			}
			r.logger.Infof("Message: %s", message)

			pr, err := r.queryLLMStream(ctx, message)
			if err != nil {
				r.logger.Error(err)
				continue
			}
			r.msg_chan <- entity.RadioMessage{
				Priority: maxPriority,
				Type:     entity.IOPIPE,
				Message:  entity.IOPipe{Pipe: pr},
				Expiry:   time.Now().Add(time.Hour),
			}
			// pr is the read end of the pipe streaming the LLM's response
			// as it's generated.

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

// queryLLMStream kicks off a streaming LLM call for the given user message
// and returns the read end of an io.Pipe that fills with the response text
// as it streams in. The producer goroutine closes the pipe (with error, if
// any) when the stream ends.
func (r *Service) queryLLMStream(ctx context.Context, userMessage string) (*io.PipeReader, error) {
	pr, pw := io.Pipe()

	go func() {
		err := llm.StreamLLM(ctx, r.systemPrompt, userMessage, r.llmModel, pw)
		_ = pw.CloseWithError(err) // nil err just closes cleanly (reader sees io.EOF)
	}()

	return pr, nil
}
