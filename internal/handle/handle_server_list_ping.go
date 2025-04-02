package handle

import (
	"github.com/blara/go-mineserver/internal/packet"
	"github.com/blara/go-mineserver/internal/state"
)

type StatusResponseStrategy struct{}

func (rs *StatusResponseStrategy) GenerateResponse(r packet.Request, s *state.Session) packet.Response {
	_ = r.(*packet.StatusRequest)
	return packet.NewStatusReponse()
}

type HandshakeResponseStrategy struct{}

func (rs *HandshakeResponseStrategy) GenerateResponse(r packet.Request, s *state.Session) packet.Response {
	req := r.(*packet.HandshakeRequest)

	s.SetState(state.SessionState(req.NextState))

	return &packet.HandshakeResponse{}
}

type PingResponseStrategy struct{}

func (rs *PingResponseStrategy) GenerateResponse(r packet.Request, s *state.Session) packet.Response {
	p := r.(*packet.PingRequest)
	resp := &packet.PingResponse{
		Timestamp: p.Timestamp,
	}

	s.Disconnect = true

	return resp
}
