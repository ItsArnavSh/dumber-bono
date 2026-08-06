package internal_types

import "context"

type STT interface {
	StartLiveSTT(ctx context.Context, chunks chan []byte) (chan []byte, error)
}

type TTS interface {
	StartStreamAudio(ctx context.Context, chunks chan []byte) (chan []byte, error)
}
