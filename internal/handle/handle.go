package handle

import (
	"bytes"
	"io"
	"net"
	"reflect"

	"github.com/blara/go-mineserver/internal/log"
	"github.com/blara/go-mineserver/internal/packet"
)

const (
	MAX_PACKET_LENGTH_BYTES uint16 = 65285
)

// HandleConnection transimits and recieves bytes on a provided connection,
// conn.
//
// The function constrains all traffic on a given client connection to a single
// abstraction.
func HandleConnection(conn net.Conn) {
	defer conn.Close()

	pkt := make([]byte, MAX_PACKET_LENGTH_BYTES) // Buffer to hold incoming data

	for {
		n, err := conn.Read(pkt)
		if err == io.EOF {
			return
		} else if err != nil {
			log.Error("Error reading from connection:", "error", err)
			return
		}

		// Create Packet from buffer
		buffer := bytes.NewBuffer(pkt[:n])
		request, err := packet.RequestFactory(buffer)
		if err != nil {
			log.Error("Error building Request", "err", err)
			break
		}
		log.Info("request created",
			"type", reflect.TypeOf(request),
			"request", log.Fmt("%+v", request),
		)

		// Build response based on which request we recieved
		var respBuffer bytes.Buffer

		responseStrategy, err := ResponseStrategyFactory(request)
		if err != nil {
			log.Error("Error", err)
			break
		}

		resp := responseStrategy.GenerateResponse(request)
		respBuffer, err = resp.Serialize()
		if err != nil {
			log.Error("Error", err)
			break
		}

		conn.Write(respBuffer.Bytes())
		log.Info("response sent",
			"type", reflect.TypeOf(resp),
			"response", log.Fmt("%+v", resp),
		)
	}
}
