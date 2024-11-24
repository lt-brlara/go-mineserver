package packet

import (
	"bytes"
	"errors"

	"github.com/blara/go-mineserver/internal/state"
)

const (
	STATUS_PACKET_ID        byte = 0x00
	PING_PACKET_ID          byte = 0x01
	CUSTOM_REPORT_PACKET_ID byte = 0x7A
)

var (
	ErrUnimplemented    = errors.New("Parser not implemented for targeted struct")
	ErrPacketNotHandled = errors.New("Packet does not have matching struct")
	ErrStateInvalid     = errors.New("Session is not in valid state for this request")
)

// Minecraft-defined hyphenated hexadecimal format (128-bit long numbers)
//
// Minecraft specifically uses "version 4, variant 1" UUIDs, where most of the
// number is randomly generated. See [RFC 4122].
//
// [RFC 4122]: https://datatracker.ietf.org/doc/html/rfc4122
type UUID string

// A Request is the generic representation of serverbound information.
type Request interface{}

// A Response is the generic interface for all clientbound information.
//
// Serialize converts the struct to byte format represented by a bytes.Buffer.
type Response interface {
	Serialize() (bytes.Buffer, error)
}

// RequestFactory returns a Request with the proper struct fields matching
// the protocol specification based on bytes read from data.
func RequestFactory(data *bytes.Buffer, session *state.Session) (Request, error) {

	length, err := readVarInt(data)
	if err != nil {
		return nil, err
	}

	packetID, err := readVarInt(data)
	if err != nil {
		return nil, err
	}
	packetIDByte := byte(packetID)

	switch packetIDByte {
	case STATUS_PACKET_ID:
		return StatusPacketFactory(packetIDByte, length, data, session)
	case PING_PACKET_ID:
		return PingPacketFactory(packetIDByte, length, data, session)
	}

	return nil, ErrPacketNotHandled
}

// StatusPacketFactory returns the correct Request based on the criteria of
// different types of status-related packets.
func StatusPacketFactory(id byte, length int32, data *bytes.Buffer, session *state.Session) (Request, error) {

	switch session.State {
	case state.StateNull:
		return NewHandshakeRequest(data)
	case state.StateStatus:
		return NewStatusRequest(data)
	case state.StateLogin:
		return NewLoginStartRequest(data)
	}

	return nil, ErrPacketNotHandled
}

// PingPacketFactory returns the correct Request based on the criteria of
// different types of ping-related packets.
func PingPacketFactory(id byte, length int32, data *bytes.Buffer, session *state.Session) (Request, error) {
	switch session.State {
	case state.StateNull:
		return nil, ErrPacketNotHandled
	case state.StateStatus:
		return NewPingRequest(data)
	case state.StateLogin:
		return NewEncryptionResponse(data)
	}

	return nil, ErrPacketNotHandled

}

func writeString(buffer *bytes.Buffer, str string) (n int, err error) {
	n1, err := writeVarInt(buffer, int32(len(str)))
	if err != nil {
		return 0, err
	}

	n2, err := buffer.WriteString(str)
	if err != nil {
		return 0, err
	}
	return int(n1) + n2, nil
}

func writeBool(buffer *bytes.Buffer, condition bool) {
	if condition {
		buffer.Write([]byte{0x01})
	} else {
		buffer.Write([]byte{0x00})
	}
}
