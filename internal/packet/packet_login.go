package packet

import (
	"bytes"

	"github.com/blara/go-mineserver/internal/crypto"
)

// A LoginStartRequest represents all fields in a serverbound intention.
//
// See the reference documentation on intention for more information:
// https://wiki.vg/Protocol#Login
type LoginStartRequest struct {
	Name       string
	PlayerUUID UUID
}

// NewLoginStartRequest returns a LoginStartRequest with all fields parsed.
func NewLoginStartRequest(data *bytes.Buffer) (*LoginStartRequest, error) {
	stringLength, err := readVarInt(data)
	if err != nil {
		return nil, err
	}
	name := make([]byte, stringLength)
	_, err = data.Read(name)

	uuid, err := readUUID(data)
	if err != nil {
		return nil, err

	}

	return &LoginStartRequest{
		Name:       string(name),
		PlayerUUID: uuid,
	}, nil
}

type EncryptionRequest struct {
	ServerID           string
	PublicKeyLength    int32
	PublicKey          []byte
	VerifyTokenLength  int32
	VerifyToken        []byte
	ShouldAuthenticate bool
}

func NewEncryptionRequest() *EncryptionRequest {

	r := &EncryptionRequest{
		ServerID:           "",
		PublicKey:          crypto.ServerKey.PublicKey,
		PublicKeyLength:    crypto.ServerKey.PublicKeyLength,
		VerifyTokenLength:  crypto.ServerKey.VerifyTokenLength,
		VerifyToken:        crypto.ServerKey.VerifyToken,
		ShouldAuthenticate: false,
	}

	return r

}

func (r *EncryptionRequest) Serialize() (bytes.Buffer, error) {
	var resp bytes.Buffer
	var buf bytes.Buffer

	writeVarInt(&buf, int32(PING_PACKET_ID))
	writeString(&buf, r.ServerID)
	writeVarInt(&buf, r.PublicKeyLength)
	buf.Write(r.PublicKey)
	writeVarInt(&buf, r.VerifyTokenLength)
	buf.Write(r.VerifyToken)
	writeBool(&buf, r.ShouldAuthenticate)

	// Calculate total length
	writeVarInt(&resp, int32(buf.Len()))

	// Copy payload buffer to response
	resp.Write(buf.Bytes())

	return resp, nil
}

// An EncryptionResponse represents all fields in a key intention.
//
// See the reference documentation on intention for more information:
// https://wiki.vg/Protocol#Login
type EncryptionResponse struct {
	SharedSecret []byte
	VerifyToken  []byte
}

// NewEncryptionResponse returns an EncryptionResponse with all fields parsed.
func NewEncryptionResponse(data *bytes.Buffer) (*EncryptionResponse, error) {
	sharedSecretLength, err := readVarInt(data)
	if err != nil {
		return nil, err
	}

	sharedSecret := make([]byte, sharedSecretLength)
	_, err = data.Read(sharedSecret)

	verifyTokenLength, err := readVarInt(data)
	if err != nil {
		return nil, err
	}

	verifyToken := make([]byte, verifyTokenLength)
	_, err = data.Read(verifyToken)
	if err != nil {
		return nil, err
	}

	return &EncryptionResponse{
		SharedSecret: sharedSecret,
		VerifyToken:  verifyToken,
	}, nil
}

// A LoginSuccess represents all fields in a game_profile resource.
//
// See the reference documentation on game_profile for more information:
// https://wiki.vg/Protocol#Login
type LoginSuccess struct {
	PlayerUUID          UUID
	Username            string
	NumberOfProperties  []Property
	StrictErrorHandling bool
}

type Property struct {
	Name      string
	Value     string
	IsSigned  bool
	Signature string
}

// NewLoginSuccess returns an EncryptionResponse with all fields parsed.
func NewLoginSuccess() *LoginSuccess {
	return &LoginSuccess{}
}

func (r *LoginSuccess) Serialize() (bytes.Buffer, error) {
	var resp bytes.Buffer
	var buf bytes.Buffer

	// Calculate total length
	writeVarInt(&resp, int32(buf.Len()))
	// Copy payload buffer to response
	resp.Write(buf.Bytes())

	return resp, nil
}
