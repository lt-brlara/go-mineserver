package packet

import (
	"bytes"
	"context"
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
func RequestFactory(ctx context.Context, data *bytes.Buffer, session *state.Session) (context.Context, Request, error) {

	length, err := readVarInt(data)
	if err != nil {
		return nil, nil, err
	}
	packetID, err := readVarInt(data)
	if err != nil {
		return nil, nil, err
	}

	ctx = NewContext(ctx, length, byte(packetID))

	switch session.State {
		case state.StateNull:
			return NullRequestFactory(ctx, data, session)
		case state.StateStatus:
			return StatusRequestFactory(ctx, data, session)
		case state.StateLogin:
			return LoginRequestFactory(ctx, data, session)
		case state.StateConfiguration:
			return ConfigurationRequestFactory(ctx, data, session)
		default:
			return ctx, nil, ErrStateNotHandled 
	}
}

func NullRequestFactory(ctx context.Context, data *bytes.Buffer, session *state.Session) (context.Context, Request, error) {

	id, ok := IdFromContext(ctx)
	if !ok { return ctx, nil, ErrPacketNotHandled }

	switch id {
		case STATUS_PACKET_ID:
			r, err := NewHandshakeRequest(data)
			return ctx, r, err
		default:
			return ctx, nil, ErrPacketNotSupported
	}
}

func StatusRequestFactory(ctx context.Context, data *bytes.Buffer, session *state.Session) (context.Context, Request, error) {

	id, ok := IdFromContext(ctx)
	if !ok { return ctx, nil, ErrPacketNotHandled }

	switch id {
		case STATUS_PACKET_ID:
			r, err := NewStatusRequest(data)
			return ctx, r, err
		case PING_PACKET_ID:
			r, err := NewPingRequest(data)
			return ctx, r, err
		default:
			return ctx, nil, ErrPacketNotSupported
	}
}

func LoginRequestFactory(ctx context.Context, data *bytes.Buffer, session *state.Session) (context.Context, Request, error) {

	id, ok := IdFromContext(ctx)
	if !ok { return ctx, nil, ErrPacketNotHandled }

	switch id {
		case STATUS_PACKET_ID:
			r, err := NewLoginStartRequest(data)
			return ctx, r, err
		case byte(0x03):
			r, err := NewLoginAcknowledgedRequest(data)
			return ctx, r, err
		default:
			return ctx, nil, ErrPacketNotSupported
	}
}

func ConfigurationRequestFactory(ctx context.Context, data *bytes.Buffer, session *state.Session) (context.Context, Request, error) {

	id, ok := IdFromContext(ctx)
	if !ok { return ctx, nil, ErrPacketNotHandled }

	switch id {
		case byte(0x07):
			r, err := NewServerboundKnownPacksRequest(data)
			return ctx, r, err
		default:
			return ctx, nil, ErrPacketNotSupported
	}
}
