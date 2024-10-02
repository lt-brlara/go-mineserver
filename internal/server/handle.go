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

	CUSTOM_REPORT_DETAILS_ID = 0x7A
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	packet := make([]byte, MAX_PACKET_LENGTH_BYTES) // Buffer to hold incoming data

	for {

		// Read data from the connection
		n, err := conn.Read(packet)
		if err != nil {
			log.Println("Error reading from connection:", err)
			return
		}

		buffer := bytes.NewBuffer(packet[:n])
		log.Printf("Received: 0x%x\n", buffer.Bytes())

		length, err := readVarInt(buffer)
		if err != nil {
			log.Println(err)
			break
		}
		log.Printf("\tLength: %d - 0x%x\n", length, byte(length))

		packetID, err := readVarInt(buffer)
		if err != nil {
			log.Println(err)
			break
		}

		log.Printf("\tPacketID: 0x%x\n", byte(packetID))

		switch byte(packetID) {
		case STATUS_PACKET_ID:
			if length > 1 {
				logStatusPacketAttributes(buffer)
			} else {
				resp, _ := handleStatusRequest()
				conn.Write(resp)
				log.Printf("Sent: 0x%x", resp)
			}
		}
	}
}

func handleStatusRequest() ([]byte, error) {
	var resp bytes.Buffer

	// Build StatusRequestResponse
	payload := &StatusRequestResponse{
		Version: Version{
			Name:     "1.21",
			Protocol: 762,
		},
		EnforcesSecureChat: false,
	}

	// byte encode JSON object

	bytePayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Could not marshal data: %v", err)
	}

	writeVarInt(&resp, int32(len(bytePayload)+1))
	writeVarInt(&resp, int32(STATUS_PACKET_ID))
	writeVarInt(&resp, int32(len(bytePayload)))
	resp.Write(bytePayload)

	return resp.Bytes(), err
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
