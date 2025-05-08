package packet

import (
	"bytes"
	"fmt"
)

type LoginStart struct {
	Name string
	PlayerUUID UUID
}

func NewLoginStartRequest(data *bytes.Buffer) (*LoginStart, error) {

	stringLength, err := readVarInt(data)
	if err != nil {
		return nil, err
	}
	name := make([]byte, stringLength)
	_, err = data.Read(name)
	uuidBytes := data.Bytes()

	_ = fmt.Sprintf( "%x-%x-%x-%x-%x",
		uuidBytes[0:4],
		uuidBytes[4:6],
		uuidBytes[6:8],
		uuidBytes[8:10],
		uuidBytes[10:],
	)

	return &LoginStart{
		Name: string(name),
		PlayerUUID: uuidBytes,
	}, nil
}

type LoginAcknowledged struct {
	Name string
	PlayerUUID UUID
}

func NewLoginAcknowledgedRequest(data *bytes.Buffer) (*LoginAcknowledged, error) {
	return &LoginAcknowledged{}, nil
}

type Property struct {
		Name string
		Value string
		Signature string
	}
type LoginSuccessResponse struct {
	UUID UUID
	Username string
	Properties []Property
}

func NewLoginSuccessResponse(r *LoginStart) (*LoginSuccessResponse) {
	return &LoginSuccessResponse{
		UUID: r.PlayerUUID,
		Username: r.Name,
	}
}

func (r *LoginSuccessResponse) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	writeVarInt(&buf, int32(LOGIN_PACKET_ID))
	_, err := buf.Write([]byte(r.UUID))
	writeVarInt(&buf, int32(len(r.Username)))
	_, err = buf.Write([]byte(r.Username))
	// TODO: Implement all parameters of login success
	writeVarInt(&buf, 0)
	writeVarInt(&buf, 0)

	var resp bytes.Buffer
	writeVarInt(&resp, int32(buf.Len()))
	resp.Write(buf.Bytes())

	return resp.Bytes(), err
}

