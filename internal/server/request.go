package server

import (
	"github.com/blara/go-mineserver/internal/packet"
	"github.com/blara/go-mineserver/internal/state"
)

type Request struct {
	Data		packet.ServerboundPacket
	Session *state.Session
}
