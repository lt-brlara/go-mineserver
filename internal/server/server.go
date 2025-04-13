package server

import (
	"bytes"
	"io"
	"net"
	"reflect"

	"github.com/blara/go-mineserver/internal/handle"
	"github.com/blara/go-mineserver/internal/log"
	"github.com/blara/go-mineserver/internal/packet"
	"github.com/blara/go-mineserver/internal/state"
)

const (
	MAX_PACKET_LENGTH_BYTES uint16 = 65285
)

var (
	serverboundChan = make(chan packet.ServerboundPacket, 100)
	tickEventChan   = make(chan Request, 100)
	clientboundChan = make(chan Response, 100)
)

type Server struct {}

func New() Server {
	return Server{}
}

// Run starts the server and creates connections to be handled downstream.
func (s *Server) Run() error {
	
	go dispatchResponse(clientboundChan)

	err := s.ListenAndServe()
	return err
}

func (s *Server) ListenAndServe() error {
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
		} 

		session := state.NewSession(conn)
		go s.Handle(session)
	}
}

func (s *Server) Handle(session *state.Session) {
	var (
		buffer *bytes.Buffer
	  tmpBuf	= make([]byte, MAX_PACKET_LENGTH_BYTES)
		readBuf = make([]byte, 0, 4096)
	)

	defer session.CloseConnection()

	for {
		n, err := session.Conn.Read(tmpBuf)
		if err == io.EOF {
			return
		} else if err != nil {
			log.Error("Error reading from connection:", "error", err)
			return
		}

		readBuf = append(readBuf, tmpBuf[:n]...)

		// Create Packet from buffer
		buffer = bytes.NewBuffer(readBuf)
		packetLength, _ := packet.GetPacketLength(buffer)

		if buffer.Len() < int(packetLength) {
			// Wait until full packet is available
			break
		}

		packetData := buffer.Next(int(packetLength))
		buf := bytes.NewBuffer(packetData)

		req := Request{
			Data: packet.Deserialize(buf, session),
			Session: session,
		}

		log.Info("request parsed",
				"type", log.Fmt("%v", reflect.TypeOf(req.Data)),
				"data", log.Fmt("%+v", req.Data),
				"session", log.Fmt("%+v", req.Session))

		if packet.IsPacketUrgent(req.Data) {
			processRequest(req)
		} else {
			serverboundChan<- req
		}
		readBuf = buffer.Bytes()
	}
}

func processRequest(r Request) {
	log.Info("request processing",
		"data", log.Fmt("%+v", r.Data),
		"session", log.Fmt("%+v", r.Session))

	rs, err := handle.ResponseStrategyFactory(r.Data)
	if err != nil {
		log.Error("error retrieving ResponseStrategy",
			"err", err,
			"response_type", reflect.TypeOf(r.Data))
			return
	}

	p, err := rs.Execute(r.Data, r.Session)
	if err != nil {
		log.Error("error executing ResponseStrategy",
			"err", err,
			"response_type", reflect.TypeOf(r.Data))
		return
	} else if p == nil {
		return
	}

	log.Info("request processed",
		"data", log.Fmt("%+v", p),
		"session", log.Fmt("%+v", r.Session))
	
	data, err := p.Serialize()
	if err != nil {
		log.Error("error serializing response", 
			"error", err,
			"response_type", reflect.TypeOf(p),
		)
		return
	}

	resp := Response{
		Data: data,
		Session: r.Session,
	}

	log.Info("response built",
		"type", reflect.TypeOf(p),
		"session", log.Fmt("%+v", r.Session),
		"resp", log.Fmt("%+v", r.Data))

	clientboundChan<- resp
}

func dispatchResponse(respChan <-chan Response) {
	for resp := range respChan {
		resp.Session.Conn.Write(resp.Data)
		log.Info("response dispatched",
			"session", log.Fmt("%+v", resp.Session),
			"data", log.Fmt("0x%x", resp.Data))
	}
}
