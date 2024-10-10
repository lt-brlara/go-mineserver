package packet

import (
	"bytes"
)

const (
	STATUS_PACKET_ID        byte = 0x00
	PING_PACKET_ID          byte = 0x01
	CUSTOM_REPORT_PACKET_ID byte = 0x7A
)

type Packet struct {
	Length   int32
	PacketID byte
	Data     []byte
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
		Data:     data.Bytes(),
	}

	return p, nil
}
