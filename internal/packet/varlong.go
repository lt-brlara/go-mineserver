package packet

import (
	"bytes"
)

// readVarLong returns the LEB-128 decoded value read from the given buffer, and
// an error if present.
func readVarLong(buffer *bytes.Buffer) (int64, error) {
	var value int64
	var position int8

	for {
		currentByte, err := buffer.ReadByte()
		if err != nil {
			return value, err
		}
		value |= (int64)(currentByte&SEGMENT_BITS) << position

		if (currentByte & CONTINUE_BIT) == 0 {
			break
		}

		position += 7

		if position >= 64 {
			return value, ErrVarLongTooBig
		}
	}

	return value, nil
}

// writeVarLong returns the number of bytes written to the given buffer, and an
// error if present.
func writeVarLong(w *bytes.Buffer, v int64) (uint8, error) {
	var n uint8
	for {
		temp := (byte)(v & 0b01111111)
		v >>= 7

		if v != 0 {
			temp |= 0b10000000
		}

		w.WriteByte(temp)
		n++

		if v == 0 {
			break
		}
	}

	return n, nil
}
