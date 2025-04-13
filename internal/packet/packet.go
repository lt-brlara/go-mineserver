package packet

import (
	"bytes"
	"errors"

	"github.com/blara/go-mineserver/internal/state"
)

const (
	STATUS_PACKET_ID        byte = 0x00
	PING_PACKET_ID          byte = 0x01
	LOGIN_PACKET_ID					byte = 0x02
	CUSTOM_REPORT_PACKET_ID byte = 0x7A
)

var (
	ErrPacketNotHandled = errors.New("Packet does not have matching struct")
	ErrPacketNotSupported = errors.New("Server does not support packet functionality")
	ErrStateNotHandled = errors.New("Packet has a state that is not handled")
)

func NullRequestFactory(id byte, data *bytes.Buffer, session *state.Session) ServerboundPacket {
	switch id {
		case byte(0x00):
			r, _ := NewHandshakeRequest(data)
			return r
		default:
			return nil
	}
}

func StatusRequestFactory(id byte, data *bytes.Buffer, session *state.Session) ServerboundPacket {

	switch id {
		case STATUS_PACKET_ID:
			r, _ := NewStatusRequest(data)
			return r
		case PING_PACKET_ID:
			r, _ := NewPingRequest(data)
			return r
		default:
			return nil
	}
}

func LoginRequestFactory(id byte, data *bytes.Buffer, session *state.Session) ServerboundPacket {
	switch id {
		case STATUS_PACKET_ID:
			r, _ := NewLoginStartRequest(data)
			return r
		case byte(0x03):
			r, _ := NewLoginAcknowledgedRequest(data)
			return r
		default:
			return nil
	}
}

func ConfigurationRequestFactory(id byte, data *bytes.Buffer, session *state.Session) ServerboundPacket {
	switch id {
		case byte(0x07):
			r, _ := NewServerboundKnownPacksRequest(data)
			return r
		case byte(0x03):
			r, _ := NewAcknowledgeFinishConfiguration(data)
			return r
		default:
			return nil
	}
}

func IsPacketUrgent(pkt ServerboundPacket) bool {
	switch pkt.(type) {
	case *HandshakeRequest,
		*LoginStartRequest,
		*LoginAcknowledgedRequest,
		*PingRequest,
		*StatusRequest,
		*ServerboundKnownPacksRequest:
		return true
	default:
		return false
	}
}

func GetPacketLength(d *bytes.Buffer) (int32, error) {
	return readVarInt(d)
}
