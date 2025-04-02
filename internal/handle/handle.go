package handle

import (
	"bytes"
	"context"
	"io"
	"reflect"

	"github.com/blara/go-mineserver/internal/log"
	"github.com/blara/go-mineserver/internal/packet"
	"github.com/blara/go-mineserver/internal/state"
)

const (
	MAX_PACKET_LENGTH_BYTES uint16 = 65285
)

// HandleConnection transimits and recieves bytes on a provided connection,
// conn.
//
// The function constrains all traffic on a given client connection to a single
// abstraction.
func HandleConnection(session *state.Session) {

	var (
		ctx context.Context
		buffer *bytes.Buffer
	)

	defer session.CloseConnection()

	pkt := make([]byte, MAX_PACKET_LENGTH_BYTES) // Buffer to hold incoming data

	for {
		n, err := session.Conn.Read(pkt)
		if err == io.EOF {
			return
		} else if err != nil {
			log.Error("Error reading from connection:", "error", err)
			return
		}

		// Create Packet from buffer
		buffer = bytes.NewBuffer(pkt[:n])
		ctx = context.Background()
		debugBytes := buffer.Bytes()

		ctx, request, err := packet.RequestFactory(ctx, buffer, session)
		id, _ := packet.IdFromContext(ctx)
		if err != nil {
			log.Error("Error building Request",
				"err", err,
				"session", log.Fmt("%+v", session),
				"packetIDByte", log.Fmt("0x%x", id),
				"data", log.Fmt("0x%x", debugBytes),
			)
			continue
		}
		log.Info("request recieved",
			"type", reflect.TypeOf(request),
			"packetID", log.Fmt("0x%x", id),
			"session", log.Fmt("%+v", session),
			"request", log.Fmt("%+v", request),
		)

		// Build response based on which request we recieved
		var respBuffer bytes.Buffer

		responseStrategy, err := ResponseStrategyFactory(request)
		if err != nil {
			log.Error("Could not retrieve ResponseStrategy", "Error", err)
			break
		}

		resp := responseStrategy.GenerateResponse(request, session)
		respBuffer, err = resp.Serialize()
		if err != nil {
			log.Error("Could not generate Response", "Error", err)
			break
		}

		session.Conn.Write(respBuffer.Bytes())
		log.Info("response transmitted",
			"type", reflect.TypeOf(resp),
			"session", log.Fmt("%+v", session),
			"response", log.Fmt("%+v", resp),
		)
		log.Debug("response bytes",
			"bytes", log.Fmt("0x%x", respBuffer.Bytes()),
		)

		if session.Disconnect {
			break
		}
	}
}
