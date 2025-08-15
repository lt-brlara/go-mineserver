package packet

import (
	"errors"
	"io"
)

const (
	SEGMENT_BITS byte = 0b01111111
	CONTINUE_BIT byte = 0b10000000
)

var (
	ErrVarIntTooBig  = errors.New("VarInt is too big")
	ErrVarLongTooBig = errors.New("VarLong is too big")
)

// readVarInt returns the LEB-128 decoded value read from the given buffer, and
// an error if present.
func ReadVarInt(r io.Reader) (int32, error) {
	var value int
	var position = 0
	var buf [1]byte

	for {
		_, err := io.ReadFull(r, buf[:])
		if err != nil {
			return 0, err
		}
		currentByte := buf[0]

		// Accumulate the result into value using lower 7 bits
		value |= int(currentByte&SEGMENT_BITS) << position

		// If MSB is not set, this is the last byte
		if (currentByte & CONTINUE_BIT) != CONTINUE_BIT {
			break
		}

		// Move to the next 7-bit block
		position += 7

		// Check for VarInt overflow
		if position >= 32 {
			return 0, ErrVarIntTooBig
		}
	}

	return int32(value), nil
}

// writeVarInt returns the number of bytes written to the given buffer, and an
// error if present.
func WriteVarInt(w io.Writer, v int32) error {
	const MAX_BYTES = 5
	var n uint8 = 0
	for {
		if n > 5 {
			return ErrVarIntTooBig
		}

		encodedBits := (byte)(v & int32(SEGMENT_BITS))
		if v & ^int32(SEGMENT_BITS) == 0 {
			w.Write([]byte{encodedBits})
			return nil
		}

		encodedBits |= CONTINUE_BIT
		w.Write([]byte{encodedBits})

		v >>= 7
		n++
	}
}
