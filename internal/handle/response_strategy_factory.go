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
	GenerateResponse(packet.Request, *state.Session) packet.Response
}

// ResponseStrategyFactory returns a ResponseStrategy interface for later use.
func ResponseStrategyFactory(req packet.Request) (ResponseStrategy, error) {
	switch req.(type) {
	case *packet.StatusRequest:
		return &StatusResponseStrategy{}, nil
	case *packet.HandshakeRequest:
		return &HandshakeResponseStrategy{}, nil
	case *packet.PingRequest:
		return &PingResponseStrategy{}, nil
	case *packet.LoginStartRequest:
		return &EncryptionRequestStrategy{}, nil
	case *packet.EncryptionResponse:
		return &LoginSuccessStrategy{}, nil
	default:
		return nil, ErrStrategyNotPresent
	}
}
