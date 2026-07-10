package audio

import (
	"context"
	"strings"
)

func (a *Audio) TTS(ctx context.Context) {
	sentencechan := make(chan string)
	go func() {
		a.Synthesize(ctx, sentencechan, a.outgoing)
	}()
	var buffer strings.Builder

	flush := func(text string) {
		text = strings.TrimSpace(text)
		if text == "" {
			return
		}

		// TODO: Send text to Piper
		sentencechan <- text
	}

	for {
		select {
		case <-ctx.Done():
			flush(buffer.String())
			return

		case chunk, ok := <-a.wordchunks:
			if !ok {
				flush(buffer.String())
				return
			}

			buffer.WriteString(" ")
			buffer.WriteString(chunk)
			s := buffer.String()

			last := strings.LastIndexAny(s, ".!?")
			if last == -1 {
				continue
			}

			// Everything up to the last '.', '!' or '?' is complete.
			flush(s[:last+1])

			buffer.Reset()

			buffer.WriteString(s[last+1:])
		}
	}
}
