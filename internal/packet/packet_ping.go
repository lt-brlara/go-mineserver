package packet

import (
	"bytes"
	"encoding/binary"
)

// A PingRequest represents all fields in a serverbound ping_request.
//
// See the reference documentation on ping_request for more information:
// https://wiki.vg/Protocol#Status
type PingRequest struct {
	Timestamp int64
}

// NewPingRequest returns a PingRequest with all fields parsed.
func NewPingRequest(data *bytes.Buffer) (*PingRequest, error) {

	// Read the integer
	var timestamp int64
	err := binary.Read(data, binary.BigEndian, &timestamp)
	if err != nil {
		return nil, err
	}

	return &PingRequest{
		Timestamp: timestamp,
	}, nil
}

type PingResponse struct {
	Timestamp int64
}

func (r *PingResponse) Serialize() ([]byte, error) {
	var resp bytes.Buffer
	var buf bytes.Buffer

	writeVarInt(&buf, int32(PING_PACKET_ID))
	err := binary.Write(&buf, binary.BigEndian, r.Timestamp)
	if err != nil {
		return nil, err
	}

	// Calculate total length
	writeVarInt(&resp, int32(buf.Len()))

	// Copy payload buffer to response
	resp.Write(buf.Bytes())

	return resp.Bytes(), nil
}
