package handle

import "github.com/blara/go-mineserver/internal/packet"

func handleClientInformation(r *Request) Result {
	req := r.Data.(*packet.ClientInformation)

	r.Client.Config.Locale = req.Locale
	r.Client.Config.ViewDistance = req.ViewDistance
	r.Client.Config.ChatMode = req.ChatMode

	return Result{}
}
