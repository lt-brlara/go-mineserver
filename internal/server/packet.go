package server

import (
	"bytes"
)

type Packet struct {
	Length int32
	PacketID byte
	Data []byte
}

func NewPacket(data *bytes.Buffer) (*Packet, error) {

	length, err := readVarInt(data)
	if err != nil {
		return &Packet{}, nil
	}

	packetID, err := readVarInt(data)
	if err != nil {
		return &Packet{}, nil
	}

	p := &Packet{
		Length: length,
		PacketID: byte(packetID),
		Data: data.Bytes(),
	}

	return p, nil
}
