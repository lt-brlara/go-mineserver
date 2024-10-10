package handle

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/blara/go-mineserver/internal/packet"
)

func handleStatusRequest() ([]byte, error) {
	// Build StatusRequestResponse
	payload := &packet.StatusRequestResponse{
		Version: packet.Version{
			Name:     "1.21.1",
			Protocol: 767,
		},
		Players: packet.Players{
			Max:    20,
			Online: 0,
		},
		Description: packet.Description{
			Text: "Test server!",
		},
		EnforcesSecureChat: false,
	}

	// byte encode JSON object
	bytePayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Could not marshal data: %v", err)
	}

	// Build payload buffer
	var buf bytes.Buffer
	packet.WriteVarInt(&buf, int32(packet.STATUS_PACKET_ID))
	packet.WriteVarInt(&buf, int32(len(bytePayload)))
	buf.Write(bytePayload)

	var resp bytes.Buffer

	// Calculate total length
	packet.WriteVarInt(&resp, int32(buf.Len()))

	// Copy payload buffer to response
	resp.Write(buf.Bytes())

	return resp.Bytes(), err
}

func handlePingRequest(req *packet.Packet) ([]byte, error) {
	var resp bytes.Buffer

	// Build payload buffer
	var buf bytes.Buffer
	packet.WriteVarInt(&buf, int32(packet.PING_PACKET_ID))
	buf.Write(req.Data)

	// Calculate total length
	packet.WriteVarInt(&resp, int32(buf.Len()))

	// Copy payload buffer to response
	resp.Write(buf.Bytes())

	return resp.Bytes(), nil
}
