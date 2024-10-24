package server

import (
	"net"

	"github.com/blara/go-mineserver/internal/handle"
	"github.com/blara/go-mineserver/internal/log"
)

func Run() error {
	listener, err := net.Listen("tcp", ":25565")
	if err != nil {
		return err
	}

	defer listener.Close()
	log.Info("Server listening on port 25565")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error("Error accepting connection", "error", err)
			continue
		} else {
			log.Info("Client connected")
		}
		go handle.HandleConnection(conn)
	}
}
