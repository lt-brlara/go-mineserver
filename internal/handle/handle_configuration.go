package handle

import (
	"github.com/blara/go-mineserver/internal/packet"
	"github.com/blara/go-mineserver/internal/state"
)


var _ ResponseStrategy = (*FinishConfigurationResponse)(nil)

func NewFinishConfigStrategy() ResponseStrategy {
	return &FinishConfigurationResponse{}
}

type FinishConfigurationResponse struct{}

func (rs *FinishConfigurationResponse) Execute(r packet.ServerboundPacket, s *state.Session) (packet.ClientboundPacket, error) {
	_ = r.(*packet.ServerboundKnownPacksRequest)
	return packet.NewClientboundKnownPacksResponse(), nil
}


var _ ResponseStrategy = (*AckFinishConfigurationResponse)(nil)

func NewAckFinishConfigurationResponse() ResponseStrategy {
	return &AckFinishConfigurationResponse{}
}

type AckFinishConfigurationResponse struct{}

func (rs *AckFinishConfigurationResponse) Execute(r packet.ServerboundPacket, s *state.Session) (packet.ClientboundPacket, error) {
	_ = r.(*packet.ServerboundKnownPacksRequest)
	return nil, nil
}
