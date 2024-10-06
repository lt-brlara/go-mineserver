package server

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
	"net"
)

const (
	STATUS_PACKET_ID byte = 0x00
	PING_PACKET_ID   byte = 0x01
	CUSTOM_REPORT_PACKET_ID byte = 0x7A
)

func handleConnection(conn net.Conn) {
	defer conn.Close()


	packet := make([]byte, MAX_PACKET_LENGTH_BYTES) // Buffer to hold incoming data

	for {
		n, err := conn.Read(packet)
		if err != nil {
			log.Println("Error reading from connection:", err)
			return
		}

		// Create Packet from buffer
		buffer := bytes.NewBuffer(packet[:n])
		request, err := NewPacket(buffer)
		if err != nil {
			log.Println(err)
			break
		}
		log.Printf("%+v", request)

		// Build response based on which request we recieved
		var resp []byte
		switch byte(request.PacketID) {
		case STATUS_PACKET_ID:
			resp, _ = handleStatusRequest()
		case PING_PACKET_ID:
			resp, _ = handlePingRequest(buffer)
		}

		// Send response to client
		conn.Write(resp)
	}
}

func handleStatusRequest() ([]byte, error) {
	// Build StatusRequestResponse
	payload := &StatusRequestResponse{
		Version: Version{
			Name:     "1.21.1",
			Protocol: 767,
		},
		Players: Players{
			Max: 20,
			Online: 0,
		},
		Description: Description{
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
	writeVarInt(&buf, int32(STATUS_PACKET_ID))
	writeVarInt(&buf, int32(len(bytePayload)))
	buf.Write(bytePayload)


	var resp bytes.Buffer

	// Calculate total length
	writeVarInt(&resp, int32(buf.Len()))

	// Copy payload buffer to response
	resp.Write(buf.Bytes())

	return resp.Bytes(), err
}

func handlePingRequest(req *bytes.Buffer) ([]byte, error) {
	var resp bytes.Buffer

	readVarLong(req)

	return resp.Bytes(), nil
}

func logStatusPacketAttributes(data *bytes.Buffer) {
	protocolVersion, err := readVarInt(data)
	if err != nil {
		log.Println(err)
	}
	log.Printf("\t\tProtocol Version: %d - 0x%x\n", protocolVersion, byte(protocolVersion))

	stringLength, err := readVarInt(data)
	if err != nil {
		log.Println(err)
	}

	serverAddr := make([]byte, stringLength)
	_, err = data.Read(serverAddr)
	if err != nil {
		log.Println("Error reading string length:", err)
	}

	log.Printf("\t\tServer Address: %s - 0x%x\n", string(serverAddr), serverAddr)

	port := make([]byte, 2)
	_, err = data.Read(port)
	if err != nil {
		log.Println("Error reading string length:", err)
	}

	log.Printf("\t\tPort: %d - 0x%x\n", binary.BigEndian.Uint16(port), port)

	nextState, err := readVarInt(data)
	if err != nil {
		log.Println(err)
	}
	log.Printf("\t\tNext State: %d - 0x%x\n", nextState, byte(nextState))

}

func logPacketAttributes(packetType string, data *bytes.Buffer) {
		log.Printf("Received: 0x%x\n", data.Bytes())

		length, err := readVarInt(data)
		if err != nil {
			log.Println(err)
		}
		log.Printf("\tLength: %d - 0x%x\n", length, byte(length))

		packetID, err := readVarInt(data)
		if err != nil {
			log.Println(err)
		}

		log.Printf("\tPacketID: 0x%x\n", byte(packetID))
}
