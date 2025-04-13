package handle

import (
	"github.com/blara/go-mineserver/internal/packet"
	"github.com/blara/go-mineserver/internal/state"
)

var _ ResponseStrategy = (*LoginStartStrategy)(nil)

func NewLoginStartStrategy() ResponseStrategy {
	return &LoginStartStrategy{}
}

type LoginStartStrategy struct{}

func (rs *LoginStartStrategy) Execute(r packet.ServerboundPacket, s *state.Session) (packet.ClientboundPacket, error) {
	req := r.(*packet.LoginStartRequest)
	return packet.NewLoginSuccessResponse(req), nil
}

var _ ResponseStrategy = (*LoginAcknowledgedStrategy)(nil)

func NewLoginAckStrategy() ResponseStrategy {
	return &LoginAcknowledgedStrategy{}
}

type LoginAcknowledgedStrategy struct{}

func (rs *LoginAcknowledgedStrategy) Execute(r packet.ServerboundPacket, s *state.Session) (packet.ClientboundPacket, error) {
	_ = r.(*packet.LoginAcknowledgedRequest)

	err := s.SetState(state.StateConfiguration)
	if err != nil {
		return nil, err
	}

	return packet.NewClientboundKnownPacksResponse(), nil
}
