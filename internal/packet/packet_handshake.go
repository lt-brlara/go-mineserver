package packet

import (
	"bytes"
	"encoding/binary"
)

// A HandshakeRequest represents all fields in a serverbound intention.
//
// See the reference documentation on intention for more information:
// https://wiki.vg/Protocol#Handshaking
type HandshakeRequest struct {
	ProtocolVersion int32  `json:"version"`
	Address         string `json:"addr"`
	Port            int32  `json:"port"`
	NextState       int32  `json:"nextState"`
}

// NewHandshakeRequest returns a HandshakeRequest with all fields parsed.
func NewHandshakeRequest(data *bytes.Buffer) (*HandshakeRequest, error) {
	protocolVersion, err := readVarInt(data)
	if err != nil {
		return nil, err
	}

	stringLength, err := readVarInt(data)
	if err != nil {
		return nil, err
	}
	serverAddr := make([]byte, stringLength)
	_, err = data.Read(serverAddr)

	port := make([]byte, 2)
	_, err = data.Read(port)
	if err != nil {
		return nil, err
	}

	nextState, err := readVarInt(data)
	if err != nil {
		return nil, err
	}

	return &HandshakeRequest{
		ProtocolVersion: protocolVersion,
		Address:         string(serverAddr),
		Port:            int32(binary.BigEndian.Uint16(port)),
		NextState:       nextState,
	}, nil
}

type HandshakeResponse struct{}

func (r *HandshakeResponse) Serialize() (bytes.Buffer, error) {
	return *bytes.NewBuffer(make([]byte, 0)), nil
}
