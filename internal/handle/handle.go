package handle

import (
	"bytes"
	"log"
	"net"

	"github.com/blara/go-mineserver/internal/packet"
)

const (
	MAX_PACKET_LENGTH_BYTES uint16 = 65285
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	pkt := make([]byte, MAX_PACKET_LENGTH_BYTES) // Buffer to hold incoming data

	for {
		n, err := conn.Read(pkt)
		if err != nil {
			log.Println("Error reading from connection:", err)
			return
		}

		// Create Packet from buffer
		buffer := bytes.NewBuffer(pkt[:n])
		log.Printf("Recieved: 0x%x", buffer)
		request, err := packet.NewPacket(buffer)
		if err != nil {
			log.Println(err)
			break
		}
		log.Printf("Request:\n\t%+v", request)

		// Build response based on which request we recieved
		var resp []byte
		switch byte(request.PacketID) {
		case packet.STATUS_PACKET_ID:
			resp, _ = handleStatusRequest()
		case packet.PING_PACKET_ID:
			resp, _ = handlePingRequest(request)
		}

		log.Printf("Sent: 0x%x", resp)
		// Send response to client
		conn.Write(resp)
	}
}
