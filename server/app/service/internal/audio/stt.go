package audio

import (
	"context"
	"fmt"
)

func (a *Audio) ProcessWAVChunks(ctx context.Context) {
	for wav := range a.incoming {
		text, err := Transcribe(ctx, wav)
		if err != nil {
			fmt.Println(err)
		}
		a.InvokeLLM(ctx, text)
	}
}
