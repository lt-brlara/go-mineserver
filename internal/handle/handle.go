package handle

import (
	"bytes"
	"io"
	"net"

	"github.com/blara/go-mineserver/internal/log"
	"github.com/blara/go-mineserver/internal/packet"
)

const (
	MAX_PACKET_LENGTH_BYTES uint16 = 65285
)

// For example..:
// Log: Info Packet recieved: PacketID = 0x00
// Log: Info Packet is Handshake: <insert fields here>

// Log: Info Packet recieved: PacketID = 0x00
// Log: Info Packet is Status Request: nil

// Log: Sending Status Reponse: <insert fields here>

// Log: Sending Pong Reponse: <insert fields here>

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	pkt := make([]byte, MAX_PACKET_LENGTH_BYTES) // Buffer to hold incoming data

	for {
		n, err := conn.Read(pkt)
		if err == io.EOF {
			return
		} else if err != nil {
			log.Error("Error reading from connection:", "error", err)
			return
		}

		// Create Packet from buffer
		buffer := bytes.NewBuffer(pkt[:n])
		request, err := packet.NewPacket(buffer)
		if err != nil {
			log.Error("Error:", err)
			break
		}

		// Build response based on which request we recieved
		var resp []byte
		switch byte(request.PacketID) {
		case packet.STATUS_PACKET_ID:
			resp, _ = handleStatusRequest(request)
		case packet.PING_PACKET_ID:
			resp, _ = handlePingRequest(request)
		}

		conn.Write(resp)

	}
}
