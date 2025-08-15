package server

import (
	"net"
	"sync"

	"github.com/blara/go-mineserver/internal/client"
)

type Server struct {
	listenAddr string
	ln         net.Listener
	quitChan   chan struct{}
	mu         sync.Mutex
	clients    []*client.Client
}

func NewServer(addr string) *Server {
	return &Server{
		listenAddr: addr,
		quitChan:   make(chan struct{}),
	}
}

func (s *Server) addClient(c *client.Client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clients = append(s.clients, c)
}

// removeClient removes a client from the server's client list.
func (s *Server) removeClient(c *client.Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, cli := range s.clients {
		if cli == c {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
			cli.Close()
			return
		}
	}
}

func (s *Server) Start() error {
	s.listen()
	defer s.ln.Close()

	go s.acceptLoop()
	go s.tick()

	<-s.quitChan

	return nil
}
