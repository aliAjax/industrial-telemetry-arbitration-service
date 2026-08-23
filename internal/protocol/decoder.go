package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
)

var (
	ErrChecksum  = errors.New("protocol checksum mismatch")
	ErrVersion   = errors.New("protocol version unsupported")
	ErrMalformed = errors.New("protocol frame malformed")
)

type Frame struct {
	Version  byte
	DeviceID uint32
	Payload  []byte
}

func Decode(data []byte) (Frame, error) {
	if len(data) < 7 {
		return Frame{}, fmt.Errorf("decode header: %w", ErrMalformed)
	}
	if data[0] != 1 {
		return Frame{}, fmt.Errorf("decode version %d: %w", data[0], ErrVersion)
	}
	payload := append([]byte(nil), data[5:len(data)-1]...)
	want := byte(0)
	for _, value := range data[:len(data)-1] {
		want ^= value
	}
	if data[len(data)-1] != want {
		return Frame{}, fmt.Errorf("decode checksum: %w", ErrChecksum)
	}
	return Frame{Version: data[0], DeviceID: binary.BigEndian.Uint32(data[1:5]), Payload: payload}, nil
}
