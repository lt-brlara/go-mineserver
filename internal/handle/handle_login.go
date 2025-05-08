package handle

import (
	"github.com/blara/go-mineserver/internal/packet"
	"github.com/blara/go-mineserver/internal/state"
)

func handleLoginStart(r *Request) Result {
	req := r.Data.(*packet.LoginStart)
	return Result{
		Response: packet.NewLoginSuccessResponse(req),
		Err:      nil,
	}
}

func handleLoginAcknowledged(r *Request) Result {
	_ = r.Data.(*packet.LoginAcknowledged)

	err := r.Client.SetState(state.Configuration)
	if err != nil {
		return Result{Response: nil, Err: err}
	}

	return Result{
		Response: packet.NewClientboundKnownPacksResponse(),
		Err:      nil,
	}
}
