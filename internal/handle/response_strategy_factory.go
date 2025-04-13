package handle

import (
	"errors"

	"github.com/blara/go-mineserver/internal/packet"
	"github.com/blara/go-mineserver/internal/state"
)

var (
	ErrStrategyNotPresent = errors.New("No strategy exists for given Request")
)

// A ResponseStrategy allows per-packet typed methods for building the set of
// response bytes to a client. ResponseStrategy should be used to interact
// with server sub-systems.
//
// GenerateReponse executes the requested action and creates the correct
// Request struct representing the return set of data to provide to a client
// connection.
type ResponseStrategy interface {
	Execute(packet.ServerboundPacket, *state.Session) (packet.ClientboundPacket, error)
}

// ResponseStrategyFactory returns a ResponseStrategy interface for later use.
func ResponseStrategyFactory(req packet.ServerboundPacket) (ResponseStrategy, error) {
	switch req.(type) {
	case *packet.HandshakeRequest:
		return NewHandshakeStrategy(), nil
	case *packet.StatusRequest:
		return NewStatusStrategy(),  nil
	case *packet.PingRequest:
		return NewPingStrategy(), nil
	case *packet.LoginStartRequest:
		return NewLoginStartStrategy(), nil
	case *packet.LoginAcknowledgedRequest:
		return NewLoginAckStrategy(), nil
	case *packet.ServerboundKnownPacksRequest:
		return NewFinishConfigStrategy(), nil
	case *packet.AcknowledgeFinishConfiguration:
		return NewAckFinishConfigurationResponse(), nil
	default:
		return nil, ErrStrategyNotPresent
	}
}
