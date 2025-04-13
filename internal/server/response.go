package server

import (
	"github.com/blara/go-mineserver/internal/log"
	"github.com/blara/go-mineserver/internal/state"
)

type Response struct {
	Data		[]byte
	Session *state.Session
}

func (r *Response) Send() {

		r.Session.Conn.Write(r.Data)
		log.Info("response transmitted",
			"response", log.Fmt("%+v", r),
		)

		log.Debug("response bytes",
			"bytes", log.Fmt("0x%x", r.Data),
		)

}

