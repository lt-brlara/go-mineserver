package packet

import (
	"bytes"
)

type ResourcePack struct {
	Namespace string
	Id string
	Version string
}

type ClientboundKnownPacksResponse struct {
	KnownPacks []ResourcePack
}

func NewClientboundKnownPacksResponse(r *LoginAcknowledgedRequest) (*ClientboundKnownPacksResponse) {
	return &ClientboundKnownPacksResponse{
		KnownPacks: []ResourcePack{
			{
				Namespace: "minecraft:core",
				Id: "1",
				Version: "1.21.1",
			},
		},
	}
}

func (r *ClientboundKnownPacksResponse) Serialize() (bytes.Buffer, error) {
	var buf bytes.Buffer

	writeVarInt(&buf, int32(byte(0x0E)))

	var prefixedArray bytes.Buffer
	writeVarInt(&prefixedArray, int32(len([]byte(r.KnownPacks[0].Namespace))))
	_, err := prefixedArray.Write([]byte(r.KnownPacks[0].Namespace))
	writeVarInt(&prefixedArray, int32(len([]byte(r.KnownPacks[0].Id))))
	_, err = prefixedArray.Write([]byte(r.KnownPacks[0].Id))
	writeVarInt(&prefixedArray, int32(len([]byte(r.KnownPacks[0].Version))))
	_, err = prefixedArray.Write([]byte(r.KnownPacks[0].Version))

	writeVarInt(&buf, int32(1))
	_, err = buf.Write(prefixedArray.Bytes())

	var resp bytes.Buffer
	writeVarInt(&resp, int32(buf.Len()))
	resp.Write(buf.Bytes())

	return resp, err
}

type ServerboundKnownPacksRequest struct {
	KnownPacks []ResourcePack
}

func NewServerboundKnownPacksRequest(data *bytes.Buffer) (*ServerboundKnownPacksRequest, error) {


	return &ServerboundKnownPacksRequest{}, nil
}
