package handle

import (
	"bytes"
	"encoding/json"

	"github.com/blara/go-mineserver/internal/log"
	"github.com/blara/go-mineserver/internal/packet"
)

func handleStatusRequest(req *packet.Packet) ([]byte, error) {

	if req.Length > 1 {
		_, err := req.ToHandshake()
		if err != nil {
			log.Error("Error converting Packet to StatusRequest", "error", err)
		}
		return nil, nil
	}

	log.Info("Packet is StatusRequest", "fields", "<nofields>")

	payload := packet.NewStatusReponse()

	// byte encode JSON object
	bytePayload, err := json.Marshal(payload)
	if err != nil {
		log.Error("Could not marshal data", "error", err)
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
	buf.Write(req.Data.Bytes())

	// Calculate total length
	packet.WriteVarInt(&resp, int32(buf.Len()))

	// Copy payload buffer to response
	resp.Write(buf.Bytes())

	return resp.Bytes(), nil
}
