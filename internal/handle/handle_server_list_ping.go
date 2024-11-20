package handle

import (
	"github.com/blara/go-mineserver/internal/packet"
)

type StatusResponseStrategy struct{}

func (s *StatusResponseStrategy) GenerateResponse(r packet.Request) packet.Response {
	_ = r.(*packet.StatusRequest)
	return packet.NewStatusReponse()
}

type HandshakeResponseStrategy struct{}

func (s *HandshakeResponseStrategy) GenerateResponse(r packet.Request) packet.Response {
	_ = r.(*packet.HandshakeRequest)
	return &packet.HandshakeResponse{}
}

type PingResponseStrategy struct{}

func (s *PingResponseStrategy) GenerateResponse(r packet.Request) packet.Response {
	p := r.(*packet.PingRequest)
	return &packet.PingResponse{
		Timestamp: p.Timestamp,
	}
}
