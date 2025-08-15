package packet

import (
	"bytes"
	"errors"
)

var (
	ErrPacketNotHandled   = errors.New("Packet does not have matching struct")
	ErrPacketNotSupported = errors.New("Server does not support packet functionality")
	ErrStateNotHandled    = errors.New("Packet has a state that is not handled")
)

func readBool(b *bytes.Buffer) (bool, error) {

	readBuffer := make([]byte, 1)
	_, err := b.Read(readBuffer)
	if err != nil {
		return false, err
	}

	if readBuffer[0] == 0x01 {
		return true, nil
	} else {
		return false, nil
	}
}

func readUint8(b *bytes.Buffer) (uint8, error) {
	readBuffer := make([]byte, 1)
	_, err := b.Read(readBuffer)
	if err != nil {
		return 0, err
	}

	return uint8(readBuffer[0]), nil
}
