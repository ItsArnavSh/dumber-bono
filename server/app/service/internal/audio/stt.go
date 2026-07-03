package audio

import (
	"context"
	"fmt"
)

func (a *Audio) ProcessWAVChunks(ctx context.Context) {
	for wav := range a.incoming {
		fmt.Println("Received the Speech Chunk")
		text, err := Transcribe(ctx, wav)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(text)
	}
}
