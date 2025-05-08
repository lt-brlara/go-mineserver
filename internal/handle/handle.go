package handle

import (
	"errors"

	"github.com/blara/go-mineserver/internal/client"
	"github.com/blara/go-mineserver/internal/packet"
	"github.com/blara/go-mineserver/internal/state"
)

const (
	MAX_PACKET_LENGTH_BYTES uint16 = 65285
)

type Request struct {
	Data    packet.Serverbound
	Client  *client.Client
	RawData []byte
}

type Result struct {
	Response packet.ClientboundPacket
	Err      error
}

type HandlerFunc func(*Request) Result

func NewRequest(data packet.Serverbound, client *client.Client, rawData []byte) Request {
	return Request{
		Data:    data,
		Client:  client,
		RawData: rawData,
	}
}

func (r *Request) Handle() Result {

	switch r.Client.State {
	case state.Null:
		return handleNull(r)
	case state.Status:
		return handleStatus(r)
	case state.Login:
		return handleLogin(r)
	case state.Configuration:
		return handleConfiguration(r)
	}

	return Result{Err: errors.New("No Result")}
}

func handleNull(r *Request) Result {
	switch r.Data.(type) {
	case *packet.HandshakeRequest:
		return handleHandshake(r)
	}

	return Result{Err: errors.New("No Result")}
}

func handleStatus(r *Request) Result {
	switch r.Data.(type) {
	case *packet.PingRequest:
		return handlePing(r)
	case *packet.StatusRequest:
		return handleStatusRequest(r)
	}

	return Result{Err: errors.New("No Result")}
}

func handleLogin(r *Request) Result {
	switch r.Data.(type) {
	case *packet.LoginStart:
		return handleLoginStart(r)
	case *packet.LoginAcknowledged:
		return handleLoginAcknowledged(r)
	}

	return Result{Err: errors.New("No Result")}
}

func handleConfiguration(r *Request) Result {
	switch r.Data.(type) {
	}

	return Result{Err: errors.New("No Result")}
}
