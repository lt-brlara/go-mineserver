package server

import (
	"fmt"
	"net"
	"reflect"

	"github.com/blara/go-mineserver/internal/client"
	"github.com/blara/go-mineserver/internal/handle"
	"github.com/blara/go-mineserver/internal/log"
	"github.com/blara/go-mineserver/internal/packet"
)

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			log.Error("error accepting connection", "error", err)
			continue
		}

		c := client.NewClient(conn)
		s.clients = append(s.clients, c)

		log.Debug("connection accepted", "client_addr", conn.RemoteAddr().String())
		go s.read(c)
	}
}

func (s *Server) read(c *client.Client) {
	defer c.Close()
	buf := make([]byte, 4096)
	for {
		n, err := c.Read(buf)
		if err != nil {
			c.HandleError(err)
			break
		}

		msg := buf[:n]

		log.Debug("message recieved from client",
			"bytes", log.Fmt("0x%x", msg),
			"state", c.State,
		)

		req := handle.NewRequest(
			packet.Parse(msg, c),
			c,
			msg,
		)

		log.Debug("\thandling packet",
			"packetType", reflect.TypeOf(req.Data),
			"state", c.State,
		)

		result := req.Handle()

		if result.Response != nil {
			resp, err := result.Response.Serialize()
			if err != nil {
				log.Error("Error serializing response", "err", err)
			}
			c.Conn.Write(resp)
		}
	}
}

func (s *Server) listen() {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		panic(err)
	}
	log.Info(fmt.Sprintf("server listening on %s", s.listenAddr))

	s.ln = ln
}

