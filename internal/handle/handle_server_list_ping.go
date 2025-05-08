package handle

import (
	"errors"

	"github.com/blara/go-mineserver/internal/packet"
)

func handleHandshake(r *Request) Result {
	var err error

	data, ok := r.Data.(*packet.HandshakeRequest)
	if !ok {
		err = errors.New("Unable to assert type PingRequest")
	}
	r.Client.SetState(data.NextState)
	return Result{
		Err: err,
	}
}

func handleStatusRequest(r *Request) Result {
	return Result{
		Response: packet.NewStatusReponse(),
		Err:      nil,
	}
}

func handlePing(r *Request) Result {
	var err error

	data, ok := r.Data.(*packet.PingRequest)
	if !ok {
		err = errors.New("Unable to assert type PingRequest")
	}
	return Result{
		Response: &packet.PingResponse{
			Timestamp: data.Timestamp,
		},
		Err: err,
	}

}
