package server

import "github.com/blara/go-mineserver/internal/packet"

type LoginStart struct {
	Name       string      `mc:"string" json:"name"`
	PlayerUUID packet.UUID `mc:"uuid" json:"playerUUID"`
}

type LoginSuccess struct {
	Packet
	UUID                packet.UUID `mc:"uuid" json:"uuid"`
	Username            string      `mc:"string" json:"username"`
	NumberOfProperties  int32       `mc:"varint" json:"numProperties"`
	Property            []Property  `mc:"array" json:"property"`
	StrictErrorHandling bool        `mc:"bool" json:"strictErrorHandling"`
}

type Property struct {
	Name      string `mc:"string" json:"name"`
	Value     string `mc:"string" json:"value"`
	IsSigned  bool   `mc:"bool" json:"isSigned"`
	Signature string `mc:"string" json:"signature,omitempty"`
}

type LoginAcknowledged struct{}

func NewLoginSuccess(req LoginStart) LoginSuccess {

	// Will likely need to do validation steps such as:
	// - Validate user info w/ Mojang as source of truth
	// - Verify user skin

	return LoginSuccess{
		Packet:              Packet{ID: 0x02},
		UUID:                req.PlayerUUID,
		Username:            req.Name,
		NumberOfProperties:  0,
		Property:            []Property{},
		StrictErrorHandling: false}
}
