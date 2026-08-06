package radio

import (
	"bufio"
	"context"
	"dubmer-bono/app/service/internal/audio/speaker"
	"io"
	"strings"
)

// speakStream reads from r incrementally, buffering text until it hits a
// sentence boundary (or EOF), then converts+plays that chunk via TTS
// immediately — rather than waiting for the whole message to arrive.
// This lets playback for streamed sources (e.g. IOPipe backed by the LLM)
// start as soon as the first sentence is ready, instead of after the
// entire response has been generated.
func (s *Service) speakStream(ctx context.Context, r io.Reader, pressure int) {
	br := bufio.NewReader(r)
	var buf strings.Builder

	flush := func() {
		text := strings.TrimSpace(buf.String())
		buf.Reset()
		if text == "" {
			return
		}

		pcm, err := s.tts.StringToPCM(ctx, text)
		if err != nil {
			s.logger.Errorf("error converting string to PCM: %v", err)
			return
		}
		s.logger.Infof("Engineer: %s \n Priority: %d", text, pressure)
		if err := speaker.PlayPCM(ctx, pcm); err != nil {
			s.logger.Errorf("error playing PCM: %v", err)
			return
		}
		s.logger.Debugf("finished playing chunk at priority %d", pressure)
	}

	for {
		ru, _, err := br.ReadRune()
		if err != nil {
			if err != io.EOF {
				s.logger.Errorf("error reading message: %v", err)
			}
			flush() // speak whatever's left in the buffer
			return
		}

		buf.WriteRune(ru)
		if strings.ContainsRune(sentenceEndings, ru) {
			flush()
		}
	}
}
