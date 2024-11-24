package handle

import (
	"github.com/blara/go-mineserver/internal/crypto"
	"github.com/blara/go-mineserver/internal/packet"
	"github.com/blara/go-mineserver/internal/state"
)

type EncryptionRequestStrategy struct{}

func (rs *EncryptionRequestStrategy) GenerateResponse(r packet.Request, s *state.Session) packet.Response {
	_ = r.(*packet.LoginStartRequest)
	return packet.NewEncryptionRequest()
}

type LoginSuccessStrategy struct{}

func (rs *LoginSuccessStrategy) GenerateResponse(r packet.Request, s *state.Session) packet.Response {
	p := r.(*packet.EncryptionResponse)

	// Perform validation of encryption

	valid, err := crypto.Validate(p.SharedSecret, p.VerifyToken, s)
	if err != nil {
	}

	// Enable encryption for session

	return packet.NewLoginSuccess()
}
