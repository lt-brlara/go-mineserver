package packet

import (
	"bytes"

	"github.com/blara/go-mineserver/internal/client"
	"github.com/blara/go-mineserver/internal/log"
	"github.com/blara/go-mineserver/internal/state"
)

type Serverbound any

func Parse(msg []byte, c *client.Client) Serverbound {
	buf := bytes.NewBuffer(msg)

	length, err := readVarInt(buf)
	if err != nil {
		log.Error("Error reading packet length", "err", err)
	}

	id, err := readVarInt(buf)
	if err != nil {
		log.Error("error reading packet ID", "msg", err)
		return nil
	}

	pkt := Packet{
		Length: length,
		ID:     byte(id),
	}

	log.Debug("\tpacket unwrapped", "pkt", log.Fmt("%+v", pkt))

	switch c.State {
	case state.Null:
		return NullRequestFactory(pkt, buf)
	case state.Status:
		return StatusRequestFactory(pkt.ID, buf)
	case state.Login:
		return LoginRequestFactory(pkt.ID, buf)
	case state.Configuration:
		return ConfigurationRequestFactory(pkt.ID, buf)
	default:
		return nil
	}
}
