package packet

import (
	"bytes"

	"github.com/blara/go-mineserver/internal/log"
)

const (
	STATUS_PACKET_ID        byte = 0x00
	PING_PACKET_ID          byte = 0x01
	CUSTOM_REPORT_PACKET_ID byte = 0x7A
)

type Packet struct {
	Length   int32
	PacketID byte
	Data     *bytes.Buffer
}

func NewPacket(data *bytes.Buffer) (*Packet, error) {

	length, err := ReadVarInt(data)
	if err != nil {
		return &Packet{}, nil
	}

	packetID, err := ReadVarInt(data)
	if err != nil {
		return &Packet{}, nil
	}

	p := &Packet{
		Length:   length,
		PacketID: byte(packetID),
		Data:     data,
	}

	// Log the packet details
	log.Info("Received Packet",
		"Length", p.Length,
		"PacketID", p.PacketID,
	)

	return p, nil
}
