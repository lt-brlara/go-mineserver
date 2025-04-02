package handle

import (
	"github.com/blara/go-mineserver/internal/packet"
	"github.com/blara/go-mineserver/internal/state"
)

type LoginStartStrategy struct{}

func (rs *LoginStartStrategy) GenerateResponse(r packet.Request, s *state.Session) packet.Response {
	req := r.(*packet.LoginStartRequest)
	return packet.NewLoginSuccessResponse(req)
}

type LoginAcknowledgedStrategy struct{}

func (rs *LoginAcknowledgedStrategy) GenerateResponse(r packet.Request, s *state.Session) packet.Response {
	req := r.(*packet.LoginAcknowledgedRequest)
	
	s.SetState(state.StateConfiguration)

	return packet.NewClientboundKnownPacksResponse(req)
}

