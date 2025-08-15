package server

import (
	"bytes"
	"fmt"
	"net"

	"github.com/blara/go-mineserver/internal/client"
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
		s.addClient(c)

		log.Debug("connection accepted", "client", c.GetAddr().String())
		go s.read(c)
	}
}

func (s *Server) read(c *client.Client) {
	defer s.removeClient(c)

	for {
		p, err := c.NextPacket()
		if err != nil {
			break
		}

		buf := bytes.NewBuffer(p)
		packetID, err := packet.ReadVarInt(buf)
		if err != nil {
			log.Error("failed to read packet ID:", err)
			continue
		}

		log.Trace("packet received", "id", packetID, "bytes", log.Fmt("0x%x", p))
		dispatch(c, packetID, buf)

		log.Trace("current client list", "clients", s.clients)
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

func dispatch(client *client.Client, id int32, data *bytes.Buffer) {
	currentState := client.GetState()

	switch {
	case currentState == 0 && id == 0x00:
		req, err := Decode[Handshake](data.Bytes())
		if err != nil {
			log.Error("failed to decode handshake", "error", err)
		}

		client.SetState(req.NextState)
		log.Debug("client state updated",
			"client", client.GetAddr().String(),
			"state", client.GetState(),
		)

	case currentState == 1 && id == 0x00:
		_, err := Decode[StatusRequest](data.Bytes())
		if err != nil {
			log.Error("failed to unmarshal handshake", "error", err)
		}

		msg := NewStatusResponse()
		err = send(client, msg)
		if err != nil {
			log.Error("failed to send message to client", "error", err)
		}

	case currentState == 1 && id == 0x01:
		req, err := Decode[PingRequest](data.Bytes())
		if err != nil {
			log.Error("failed to decode ping request", "error", err)
		}

		msg := NewPingResponse(req.Timestamp)
		send(client, msg)

	case currentState == 2 && id == 0x00:
		req, err := Decode[LoginStart](data.Bytes())
		if err != nil {
			log.Error("failed to decode login start", "error", err)
		}

		log.Info("player connecting",
			"username", req.Name,
			"uuid", req.PlayerUUID.FormattedString(),
		)

		// Will likely need to do validation steps such as:
		// - Validate user info w/ Mojang as source of truth
		// - Verify user skin

		client.SetUsername(req.Name)
		client.SetUUID(req.PlayerUUID)

		msg := NewLoginSuccess(req)
		send(client, msg)

	case currentState == 2 && id == 0x03:
		_, err := Decode[LoginAcknowledged](data.Bytes())
		if err != nil {
			log.Error("failed to decode login acknowledged", "error", err)
		}

		log.Info("player logged in",
			"username", client.GetUsername(),
			"uuid", client.GetUUID().FormattedString(),
		)

		client.SetState(3)

	default:
		log.Error("unknown packet ID",
			"id", id,
			"clientState", client.GetState(),
			"bytes", log.Fmt("%02x%x", byte(id), data),
		)
	}
}

func send(c *client.Client, msg any) error {
	pkt, err := Encode(msg)
	if err != nil {
		return err
	}

	log.Trace("sending message", "client_addr", c.GetAddr(), "bytes", log.Fmt("%x", pkt))

	_, err = c.Write(pkt)
	return err
}
