package packet

import (
	"bytes"

	"github.com/blara/go-mineserver/internal/log"
)

type ClientInformation struct {
	Locale                string
	ViewDistance          int8
	ChatMode              int32
	ChatColors            bool
	DisplayedSkinParts    uint8
	MainHand              int32
	TextFilteringEnabled  bool
	ServerListingsEnabled bool
}

func NewClientInformation(data *bytes.Buffer) (*ClientInformation, error) {
	var info ClientInformation
	var readBuffer []byte

	localeLength, err := readVarInt(data)
	if err != nil {
		return nil, err
	}

	readBuffer = make([]byte, localeLength)
	data.Read(readBuffer)
	info.Locale = string(readBuffer)

	readBuffer = make([]byte, 1)
	data.Read(readBuffer)
	info.ViewDistance = int8(readBuffer[0])

	chatMode, err := readVarInt(data)
	if err != nil {
		return &info, err
	}
	info.ChatMode = chatMode

	chatColorsEnabled, err := readBool(data)
	if err != nil {
		return &info, err
	}
	info.ChatColors = chatColorsEnabled

	displayedSkinParts, err := readUint8(data)
	if err != nil {
		return &info, err
	}
	info.DisplayedSkinParts = displayedSkinParts

	mainHand, err := readVarInt(data)
	if err != nil {
		return &info, err
	}
	info.MainHand = mainHand

	filteringEnabled, err := readBool(data)
	if err != nil {
		return &info, err
	}
	info.TextFilteringEnabled = filteringEnabled

	listingsEnabled, err := readBool(data)
	if err != nil {
		return &info, err
	}
	info.ServerListingsEnabled = listingsEnabled

	log.Debug("\tparsed client info", "packet", log.Fmt("%+v", info))
	return &info, nil
}

type ResourcePack struct {
	Namespace string
	Id        string
	Version   string
}

type ClientboundKnownPacksResponse struct {
	KnownPacks []ResourcePack
}

func NewClientboundKnownPacksResponse() *ClientboundKnownPacksResponse {
	return &ClientboundKnownPacksResponse{
		KnownPacks: []ResourcePack{
			{
				Namespace: "minecraft:core",
				Id:        "1",
				Version:   "1.21.1",
			},
		},
	}
}

func (r *ClientboundKnownPacksResponse) Serialize() ([]byte, error) {
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

	return resp.Bytes(), err
}

type ServerboundKnownPacksRequest struct {
	KnownPacks []ResourcePack
}

func NewServerboundKnownPacksRequest(data *bytes.Buffer) (*ServerboundKnownPacksRequest, error) {
	return &ServerboundKnownPacksRequest{}, nil
}

type AcknowledgeFinishConfiguration struct{}

func NewAcknowledgeFinishConfiguration(data *bytes.Buffer) (*AcknowledgeFinishConfiguration, error) {
	return &AcknowledgeFinishConfiguration{}, nil
}
