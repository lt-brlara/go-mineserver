package handle

const (
	MAX_PACKET_LENGTH_BYTES uint16 = 65285
)

// HandleConnection transimits and recieves bytes on a provided connection,
// conn.
//
// The function constrains all traffic on a given client connection to a single
// abstraction.

