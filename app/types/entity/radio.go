package entity

import (
	"io"
	"time"
)

type PilotPressurePhysicalFactors struct {
	// Turn Indicators
	GForceLat   float32
	GFroceLon   float32
	Steer       float32
	Brake       float32
	LapDistance float32
}

type MessageType int
type RadioPayload interface {
	isRadioPayload()
}

type DirectMessage struct {
	Text string
}

func (DirectMessage) isRadioPayload() {}

type FunctionMessage struct {
	Fn func() string
}

func (FunctionMessage) isRadioPayload() {}

type IOPipe struct {
	Pipe *io.PipeReader
}

func (IOPipe) isRadioPayload() {}

const (
	DIRECT MessageType = iota
	FUNCTION
	IOPIPE
)

type RadioMessage struct {
	Priority int
	Message  RadioPayload
	Type     MessageType
	Expiry   time.Time
}

func (r RadioMessage) GetExpiry() time.Time {
	return r.Expiry
}
