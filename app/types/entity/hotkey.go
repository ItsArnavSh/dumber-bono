package entity

type HotKeyEvent int

const (
	RADIO_PRESS HotKeyEvent = iota
	RADIO_RELEASE
	COPY_AFFIRMATION
	MUTE_TOGGLE //Does not mute level 5 inst
)
