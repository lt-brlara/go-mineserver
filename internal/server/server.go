package server

import (
	"log"
	"net"

	"github.com/blara/go-mineserver/internal/handle"
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
		go handle.HandleConnection(conn)
	}
}
