package server

import (
	"log"
	"net"
)

const (
	MAX_PACKET_LENGTH_BYTES uint16 = 65285
)

func Run() error {
	listener, err := net.Listen("tcp", ":25565")
	if err != nil {
		return err
	}

	defer listener.Close()
	log.Println("Server listening on port 25565")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		} else {
			log.Println("Client connected")
		}
		go handleConnection(conn)
	}
}

type Player struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type Version struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type Players struct {
	Max    int      `json:"max"`
	Online int      `json:"online"`
	//Sample []Player `json:"sample"`
}

type Description struct {
	Text string `json:"text"`
}

type StatusRequestResponse struct {
	Version            `json:"version"`
	Players            `json:"players"`
	Description        `json:"description"`
	Favicon            string `json:"favicon"`
	EnforcesSecureChat bool   `json:"enforcesSecureChat"`
}
