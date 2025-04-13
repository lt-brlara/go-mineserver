package packet

import (
	"bytes"

	"github.com/blara/go-mineserver/internal/log"
	"github.com/blara/go-mineserver/internal/state"
)

type ServerboundPacket any

func Deserialize(buf *bytes.Buffer, s *state.Session) ServerboundPacket {
	packetID, err := readVarInt(buf)
	if err != nil {
		log.Error("error reading packet ID", "msg", err)
		return nil
	}

	switch s.State {
		case state.StateNull:
			return NullRequestFactory(byte(packetID), buf, s)
		case state.StateStatus:
			return StatusRequestFactory(byte(packetID), buf, s)
		case state.StateLogin:
			return LoginRequestFactory(byte(packetID), buf, s)
		case state.StateConfiguration:
			return ConfigurationRequestFactory(byte(packetID), buf, s)
		default:
			return nil
	}
}
