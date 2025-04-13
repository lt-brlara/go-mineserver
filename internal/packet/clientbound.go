package packet

type ClientboundPacket interface {
	Serialize() ([]byte, error)
}
