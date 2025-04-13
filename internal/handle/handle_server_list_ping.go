package handle

import (
	"github.com/blara/go-mineserver/internal/packet"
	"github.com/blara/go-mineserver/internal/state"
)


var _ ResponseStrategy = (*StatusStrategy)(nil)

func NewStatusStrategy() ResponseStrategy {
    return &StatusStrategy{}
}

type StatusStrategy struct{}

func (rs *StatusStrategy) Execute(r packet.ServerboundPacket, s *state.Session) (packet.ClientboundPacket, error) {
	_ = r.(*packet.StatusRequest)
	return packet.NewStatusReponse(), nil
}


var _ ResponseStrategy = (*HandshakeResponseStrategy)(nil)
func NewHandshakeStrategy() ResponseStrategy {
    return &HandshakeResponseStrategy{}
}

type HandshakeResponseStrategy struct{}

func (rs *HandshakeResponseStrategy) Execute(r packet.ServerboundPacket, s *state.Session) (packet.ClientboundPacket, error) {
	req := r.(*packet.HandshakeRequest)

	s.SetState(state.SessionState(req.NextState))

	return nil, nil
}


var _ ResponseStrategy = (*PingStrategy)(nil)

func NewPingStrategy() ResponseStrategy {
    return &PingStrategy{}
}

type PingStrategy struct{}

func (rs *PingStrategy) Execute(r packet.ServerboundPacket, s *state.Session) (packet.ClientboundPacket, error) {
	p := r.(*packet.PingRequest)
	resp := &packet.PingResponse{
		Timestamp: p.Timestamp,
	}

	s.Disconnect = true

	return resp, nil
}
