package packet

import (
	"bytes"
	"errors"

	"github.com/blara/go-mineserver/internal/log"
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

type Packet struct {
	Length int32
	ID byte
}

func readBool(b *bytes.Buffer) (bool, error) {

	readBuffer := make([]byte, 1)
	_, err := b.Read(readBuffer); if err != nil {
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
	_, err := b.Read(readBuffer); if err != nil {
		return 0, err
	}

	return uint8(readBuffer[0]), nil
}

func NullRequestFactory(pkt Packet, data *bytes.Buffer) Serverbound {
	switch pkt.ID {
	case byte(0x00):
		r, _ := NewHandshakeRequest(data)
		return r
	default:
		return nil
	}
}

func StatusRequestFactory(id byte, data *bytes.Buffer) Serverbound {

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

func LoginRequestFactory(id byte, data *bytes.Buffer) Serverbound {
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

func ConfigurationRequestFactory(id byte, data *bytes.Buffer) Serverbound {
	switch id {
	case byte(0x00):
		r, err := NewClientInformation(data); if err != nil {
			log.Error("Error parsing packet", "err", err, "data", log.Fmt("%+v", r))
		}
		return r
	case byte(0x03):
		r, err := NewAcknowledgeFinishConfiguration(data); if err != nil {
			log.Error("Error parsing packet", "err", err, "data", log.Fmt("%+v", r))
		}
		return r
	case byte(0x07):
		r, err := NewServerboundKnownPacksRequest(data); if err != nil {
			log.Error("Error parsing packet", "err", err, "data", log.Fmt("%+v", r))
		}
		return r
	default:
	return nil
	}
}
