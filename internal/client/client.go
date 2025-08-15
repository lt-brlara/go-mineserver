package client

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/blara/go-mineserver/internal/log"
	"github.com/blara/go-mineserver/internal/packet"
)

const MAX_PACKET_LENGTH = 2097151 // 2^21 - 1

type Client struct {
	conn     net.Conn
	state    int32
	mu       sync.Mutex
	username string
	uuid     packet.UUID
}

func NewClient(c net.Conn) *Client {
	return &Client{
		conn:  c,
		state: 0,
	}
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Read(b []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.conn.Read(b)
}

func (c *Client) Write(b []byte) (n int, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.conn.Write(b)
}

func (c *Client) GetState() int32 {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.state
}

func (c *Client) SetState(s int32) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.state = s
}

func (c *Client) HandleError(err error) {
	if err == io.EOF {
		log.Info("closing connection", "client", c.conn.RemoteAddr(), "state", c.state)
	} else {
		log.Error("error reading from connection", "err", err, "state", c.state)
	}
	c.Close()
}

func (c *Client) NextPacket() ([]byte, error) {
	length, err := packet.ReadVarInt(c)
	if err != nil {
		return nil, err
	}

	if length <= 0 || length > MAX_PACKET_LENGTH {
		return nil, fmt.Errorf("invalid packet length: %d", length)
	}

	packetBytes := make([]byte, length)
	_, err = io.ReadFull(c, packetBytes)
	if err != nil {
		return nil, err
	}

	return packetBytes, nil // This is [PacketID+Data]
}

func (c *Client) GetAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Client) GetUsername() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.username
}

func (c *Client) SetUsername(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.username = name
}

func (c *Client) GetUUID() packet.UUID {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.uuid
}

func (c *Client) SetUUID(uuid packet.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.uuid = uuid
}
