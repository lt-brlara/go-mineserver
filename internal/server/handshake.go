package server

type Handshake struct {
	ProtocolVersion int32  `mc:"varint" json:"version"`
	Address         string `mc:"string" json:"address"`
	Port            uint16 `mc:"unsignedshort" json:"port"`
	NextState       int32  `mc:"varint" json:"nextState"`
}
