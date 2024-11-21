package state

import (
	"net"
)

// A Session encapsulates all relevant data for a specific client connection.
type Session struct {
	State SessionState
	Conn  net.Conn
}

type SessionState uint8

// All possible client states
const (
	StateNull SessionState = iota
	StateStatus
	StateLogin
	StateTransfer
)

// NewSession returns a pointer to Session encapsulating all relevant state data.
func NewSession(conn net.Conn) *Session {
	return &Session{
		State: StateNull,
		Conn:  conn,
	}
}

// SetState un-safely redefines the clients state based solely on the clients
// instruction.
func (s *Session) SetState(desiredState SessionState) error {
	// TODO: implement validation
	s.State = desiredState
	return nil
}
