package client

import (
	"io"
	"net"
	"sync"

	"github.com/blara/go-mineserver/internal/log"
	"github.com/blara/go-mineserver/internal/state"
)

type Client struct {
	Conn   net.Conn
	State  state.State
	mu     sync.Mutex
	Config ClientConfig
}

type ClientConfig struct {
	Locale                string
	ViewDistance          int8
	ChatMode              int32
	ChatColors            bool
	DisplayedSkinParts    uint8
	MainHand              int32
	TextFilteringEnabled  bool
	ServerListingsEnabled bool
}

func NewClient(c net.Conn) *Client {
	return &Client{
		Conn:  c,
		State: state.Null,
	}
}

func (c *Client) Close() {
	c.Conn.Close()
}

func (c *Client) Read(b []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.Conn.Read(b)
}

func (c *Client) Write(b []byte) (n int, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.Conn.Write(b)
}

func (c *Client) SetState(s state.State) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.State = s
	return nil
}

func (c *Client) HandleError(err error) {
	if err == io.EOF {
		log.Info("closing connection", "client", c.Conn.RemoteAddr(), "state", c.State)
	} else {
		log.Error("error reading from connection", "err", err, "state", c.State)
	}
	c.Close()
}
