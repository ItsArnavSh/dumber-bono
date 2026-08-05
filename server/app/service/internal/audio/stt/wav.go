package audio

import (
	"bytes"
	"encoding/binary"
)

func PCM16ToWAV(samples []int16, sampleRate, channels int) ([]byte, error) {
	const bitsPerSample = 16

	dataSize := len(samples) * 2
	byteRate := sampleRate * channels * bitsPerSample / 8
	blockAlign := channels * bitsPerSample / 8

	buf := new(bytes.Buffer)

	// RIFF header
	buf.WriteString("RIFF")
	if err := binary.Write(buf, binary.LittleEndian, uint32(36+dataSize)); err != nil {
		return nil, err
	}
	buf.WriteString("WAVE")

	// fmt chunk
	buf.WriteString("fmt ")
	if err := binary.Write(buf, binary.LittleEndian, uint32(16)); err != nil { // PCM header size
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, uint16(1)); err != nil { // PCM format
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, uint16(channels)); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, uint32(sampleRate)); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, uint32(byteRate)); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, uint16(blockAlign)); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, uint16(bitsPerSample)); err != nil {
		return nil, err
	}

	// data chunk
	buf.WriteString("data")
	if err := binary.Write(buf, binary.LittleEndian, uint32(dataSize)); err != nil {
		return nil, err
	}

	// PCM samples
	if err := binary.Write(buf, binary.LittleEndian, samples); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
