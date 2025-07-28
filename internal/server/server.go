package server

import (
	"net"

	"github.com/blara/go-mineserver/internal/client"
	"github.com/blara/go-mineserver/internal/handle"
)

type Server struct {
	listenAddr string
	ln         net.Listener
	quitChan   chan struct{}
	resultChan chan handle.Result
	clients    []*client.Client
	eventQueue *Queue
}

func NewServer(addr string) *Server {
	return &Server{
		listenAddr: addr,
		quitChan:   make(chan struct{}),
		resultChan: make(chan handle.Result, 100),
		eventQueue: NewEventQueue(),
	}
}

func (s *Server) Start() error {

	s.listen()
	defer s.ln.Close()

	go s.acceptLoop()
	go s.tick()

	<-s.quitChan
	close(s.resultChan)

	return nil
}
