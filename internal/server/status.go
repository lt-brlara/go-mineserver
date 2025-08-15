package server

import (
	"encoding/json"
)

type StatusRequest struct{}

type StatusResponse struct {
	Packet
	JsonResponse []byte `mc:"string"`
}

type Status struct {
	Version            StatusVersion     `json:"version"`
	Players            StatusPlayers     `json:"players"`
	Description        StatusDescription `json:"description"`
	Favicon            string            `json:"favicon"`
	EnforcesSecureChat bool              `json:"enforcesSecureChat"`
}

type StatusVersion struct {
	Name     string `json:"name"`
	Protocol int32  `json:"protocol"`
}

type StatusPlayers struct {
	Max    int32    `json:"max"`
	Online int32    `json:"online"`
	Sample struct{} `json:"sample"`
}

type StatusDescription struct {
	Text string `json:"text"`
}

type PingRequest struct {
	Timestamp int64 `mc:"long" json:"timestamp"`
}

type PingResponse struct {
	Packet
	Timestamp int64 `mc:"long" json:"timestamp"`
}

func NewStatusResponse() StatusResponse {
	status := NewStatus()
	b, err := json.Marshal(&status)
	if err != nil {
		panic(err)
	}

	return StatusResponse{
		Packet:       Packet{ID: 0x00},
		JsonResponse: b,
	}
}

func NewStatus() Status {
	return Status{
		Version: StatusVersion{
			Name:     "1.21.1",
			Protocol: 767},
		Players: StatusPlayers{
			Max:    20,
			Online: 0,
			Sample: struct{}{}},
		Description: StatusDescription{
			Text: "This is a description"},
		Favicon:            "",
		EnforcesSecureChat: false}
}

func NewPingResponse(timestamp int64) PingResponse {
	return PingResponse{
		Packet:    Packet{ID: 0x01},
		Timestamp: timestamp}
}
