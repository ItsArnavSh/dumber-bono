package parsers

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type PacketHeader struct {
	PacketFormat            uint16
	GameYear                uint8
	GameMajorVersion        uint8
	GameMinorVersion        uint8
	PacketVersion           uint8
	PacketID                uint8
	SessionUID              uint64
	SessionTime             float32
	FrameIdentifier         uint32
	OverallFrameIdentifier  uint32
	PlayerCarIndex          uint8
	SecondaryPlayerCarIndex uint8
}

var HeaderSize = binary.Size(PacketHeader{})

func ParseHeader(data []byte) (*PacketHeader, error) {
	if len(data) < HeaderSize {
		return nil, fmt.Errorf("parsers: payload too short (%d bytes, need %d)", len(data), HeaderSize)
	}
	var h PacketHeader
	if err := binary.Read(bytes.NewReader(data[:HeaderSize]), binary.LittleEndian, &h); err != nil {
		return nil, err
	}
	return &h, nil
}
