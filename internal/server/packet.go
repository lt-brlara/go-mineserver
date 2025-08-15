package server

type Packet struct {
	ID int32 `mc:"varint" json:"id"`
}
