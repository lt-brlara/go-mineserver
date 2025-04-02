package state

import (
	"context"
	"net"

	"github.com/blara/go-mineserver/internal/log"
)

// A Session encapsulates all relevant data for a specific client connection.
type Session struct {
	State SessionState
	Conn  net.Conn
	Disconnect bool
	Ctx context.Context
}

type SessionState uint8

// NewSession returns a pointer to Session encapsulating all relevant state data.
func NewSession(conn net.Conn) *Session {
	s := &Session{
		State: StateNull,
		Conn:  conn,
	}

	log.Info("Client connected",
		"session", log.Fmt("%+v", s),
	)
	return s
}

func (s *Session) CloseConnection() {
	err := s.Conn.Close()
	if err != nil {
		log.Error("Client had ungraceful disconnection", 
			"err", err,
		)
	}
	log.Info("Client disconnected",
		"session", log.Fmt("%+v", s),
	)
}

// SetState un-safely redefines the clients state based solely on the clients
// instruction.
func (s *Session) SetState(desiredState SessionState) error {
	// TODO: implement validation
	s.State = desiredState
	return nil
}
