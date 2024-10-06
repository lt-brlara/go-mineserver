package server

import (
	"bytes"
	"testing"
)

func TestReadVarInt(t *testing.T) {

	readTest := bytes.NewBuffer([]byte{
		0x00,
		0x01,
		0x02,
		0x03,
		0x80, 0x01,
		0xfe, 0x01,
		0xff, 0x01,
		0xff, 0xff, 0x7f,
		0xff, 0xff, 0xff, 0xff, 0x07,
		0xff, 0xff, 0xff, 0xff, 0x0f,
		0x80, 0x80, 0x80, 0x80, 0x08,
	})

	writeTest := []int32{
		0,
		1,
		2,
		3,
		128,
		254,
		255,
		2097151,
		2147483647,
		-1,
		-2147483648,
	}

	for _, w := range writeTest {
		result, err := readVarInt(readTest)
		if err != nil {
			t.Fatal(err)
		}

		if result != w {
			t.Fatalf("VarInt mismatch. Expected: %d, Got: %d", w, result)
		}
	}
}

func TestWriteVarInt(t *testing.T) {
	tests := []struct {
		input    int32
		expected []byte
	}{
		{1, []byte{0x01}},
		{2, []byte{0x02}},
		{127, []byte{0x7F}},
		{128, []byte{0x80, 0x01}},
		{300, []byte{0xAC, 0x02}},
	}

	for _, tt := range tests {
		t.Run("Encoding int32", func(t *testing.T) {

			var buf bytes.Buffer
			_, err := writeVarInt(&buf, tt.input)
			if err != nil {
				t.Fatalf("Failed to write VarInt: %v", err)
			}
			if !bytes.Equal(buf.Bytes(), tt.expected) {
				t.Errorf("Expected: %x, Got: %x", tt.expected, buf.Bytes())
			}
		})
	}
}

func TestReadVarLong(t *testing.T) {

	readTest := bytes.NewBuffer([]byte{
		0x00,
		0x01,
		0x02,
		0x03,
		0x80, 0x01,
		0xfe, 0x01,
		0xff, 0x01,
		0xff, 0xff, 0x7f,
		0xff, 0xff, 0xff, 0xff, 0x07,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 
		0x80, 0x80, 0x80, 0x80, 0xf8, 0xff, 0xff, 0xff, 0xff, 0x01,
		0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x01,
	})

	writeTest := []int64{
		0,
		1,
		2,
		3,
		128,
		254,
		255,
		2097151,
		2147483647,
		9223372036854775807,
		-1,
		-2147483648,
		-9223372036854775808,
	}

	for _, w := range writeTest {
		result, err := readVarLong(readTest)
		if err != nil {
			t.Fatal(err)
		}

		if result != w {
			t.Fatalf("VarInt mismatch. Expected: %d, Got: %d", w, result)
		}
	}
}
