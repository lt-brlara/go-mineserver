package packet

import (
	"encoding/binary"

	"github.com/blara/go-mineserver/internal/log"
)

type Player struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type Version struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type Players struct {
	Max    int `json:"max"`
	Online int `json:"online"`
	//Sample []Player `json:"sample"`
}

type Description struct {
	Text string `json:"text"`
}

type StatusResponse struct {
	Version            `json:"version"`
	Players            `json:"players"`
	Description        `json:"description"`
	Favicon            string `json:"favicon"`
	EnforcesSecureChat bool   `json:"enforcesSecureChat"`
}

type Handshake struct {
	ProtocolVersion int32  `json:"version"`
	Address         string `json:"addr"`
	Port            int32  `json:"port"`
	NextState       int32  `json:"nextState"`
}

func (p *Packet) ToHandshake() (req Handshake, err error) {
	protocolVersion, err := ReadVarInt(p.Data)
	if err != nil {
		return
	}
	req.ProtocolVersion = protocolVersion

	stringLength, err := ReadVarInt(p.Data)
	if err != nil {
		return
	}
	serverAddr := make([]byte, stringLength)
	_, err = p.Data.Read(serverAddr)
	req.Address = string(serverAddr)

	port := make([]byte, 2)
	_, err = p.Data.Read(port)
	if err != nil {
		return
	}
	req.Port = int32(binary.BigEndian.Uint16(port))

	nextState, err := ReadVarInt(p.Data)
	if err != nil {
		return
	}
	req.NextState = nextState

	log.Info("Packet is Handshake",
		"protocolVersion", req.ProtocolVersion,
		"address", req.Address,
		"port", req.Port,
		"nextState", req.NextState,
	)
	return
}

func NewStatusReponse() *StatusResponse {
	resp := &StatusResponse{
		Version: Version{
			Name:     "1.21.1",
			Protocol: 767,
		},
		Players: Players{
			Max:    20,
			Online: 0,
		},
		Description: Description{
			Text: "Test server!",
		},
		EnforcesSecureChat: false,
	}

	return resp
}
