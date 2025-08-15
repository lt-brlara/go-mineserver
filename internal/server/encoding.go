package server

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"

	"github.com/blara/go-mineserver/internal/packet"
)

func Decode[T any](b []byte) (T, error) { var v T; return v, Unmarshal(b, &v) }

func Encode(v any) ([]byte, error) {
	payload, err := Marshal(v)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := packet.WriteVarInt(&buf, int32(len(payload))); err != nil {
		return nil, err
	}

	buf.Write(payload)
	return buf.Bytes(), nil
}

func Unmarshal(data []byte, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("must pass a non-nil pointer")
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return errors.New("must pass pointer to struct")
	}

	buf := bytes.NewBuffer(data)

	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		value := rv.Field(i)

		// Recurse into embedded structs
		if field.Anonymous {
			if err := Unmarshal(buf.Bytes()[len(data)-buf.Len():], value.Addr().Interface()); err != nil {
				return err
			}
			continue
		}

		tag := field.Tag.Get("mc")

		switch tag {
		case "varint":
			v, err := packet.ReadVarInt(buf)
			if err != nil {
				return err
			}
			value.SetInt(int64(v))

		case "string":
			length, err := packet.ReadVarInt(buf)
			if err != nil {
				return err
			}
			v := make([]byte, length)
			_, err = buf.Read(v)
			value.SetString(string(v))

		case "unsignedshort":
			v := make([]byte, 2)
			_, err := buf.Read(v)
			if err != nil {
				return err
			}
			value.SetUint(uint64(binary.BigEndian.Uint16(v)))

		case "long":
			// Read the integer
			var v int64
			err := binary.Read(buf, binary.BigEndian, &v)
			if err != nil {
				return err
			}
			value.SetInt(v)

		case "uuid":
			uuid := make([]byte, 16)
			_, err := buf.Read(uuid)
			if err != nil {
				return err
			}

			value.SetBytes(uuid)

		default:
			return fmt.Errorf("unsupported tag '%s'", tag)
		}
	}

	return nil
}

// Marshal inspects struct fields tagged with `mc` and writes them in order.
// Supported tags: "varint", "string", "unsignedshort". Anonymous embedded
// structs are recursed first (so Packet.ID is output before other fields).
func Marshal(v any) ([]byte, error) {
	rv := reflect.ValueOf(v)

	if !rv.IsValid() {
		return nil, errors.New("must pass a non-nil struct or pointer to struct")
	}

	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return nil, errors.New("must pass a non-nil struct or pointer to struct")
		}
		rv = rv.Elem()
	}

	// ➜ Make it addressable so rv.Field(i).Addr() is valid
	if !rv.CanAddr() {
		av := reflect.New(rv.Type()).Elem()
		av.Set(rv)
		rv = av
	}

	var buf bytes.Buffer
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fv := rv.Field(i)

		// Recurse into anonymous embedded structs first
		if field.Anonymous {
			b, err := Marshal(fv.Addr().Interface())
			if err != nil {
				return nil, err
			}
			buf.Write(b)
			continue
		}

		tag := field.Tag.Get("mc")
		switch tag {
		case "varint":
			if err := packet.WriteVarInt(&buf, int32(fv.Int())); err != nil {
				return nil, err
			}

		case "string":
			if err := encodeString(fv, field, &buf); err != nil {
				return nil, err
			}

		case "unsignedshort":
			if err := binary.Write(&buf, binary.BigEndian, uint16(fv.Uint())); err != nil {
				return nil, err
			}

		case "long":
			err := binary.Write(&buf, binary.BigEndian, fv.Int())
			if err != nil {
				return nil, err
			}

		case "bool":
			if fv.Kind() != reflect.Bool {
				return nil, fmt.Errorf("field %s: bool tag on %s", field.Name, fv.Kind())
			}
			if fv.Bool() {
				buf.WriteByte(1)
			} else {
				buf.WriteByte(0)
			}

		case "uuid":
			v, ok := fv.Interface().(packet.UUID)
			if !ok {
				return nil, fmt.Errorf("field %s: uuid tag expects packet.UUID", field.Name)
			}

			_, err := buf.Write(v)
			if err != nil {
				return nil, err
			}

		case "array":
			if fv.Kind() != reflect.Slice && fv.Kind() != reflect.Array {
				return nil, fmt.Errorf("field %s: array tag requires slice/array", field.Name)
			}

			encodeArray(fv, field, &buf)

		case "":
			// no mc tag: skip
		default:
			return nil, fmt.Errorf("unsupported mc tag %q on field %s", tag, field.Name)
		}
	}

	return buf.Bytes(), nil
}

func encodeArray(v reflect.Value, field reflect.StructField, buf *bytes.Buffer) error {
	for j := 0; j < v.Len(); j++ {
		element := v.Index(j)
		switch {
		case element.Kind() == reflect.Struct:
			b, err := Marshal(element.Addr().Interface())
			if err != nil {
				return fmt.Errorf("%s[%d]: %w", field.Name, j, err)
			}
			buf.Write(b)
		case element.Kind() == reflect.Ptr && !element.IsNil() && element.Elem().Kind() == reflect.Struct:
			b, err := Marshal(element.Interface())
			if err != nil {
				return fmt.Errorf("%s[%d]: %w", field.Name, j, err)
			}
			buf.Write(b)
		default:
			return fmt.Errorf("field %s: array elements must be struct or *struct, got %s", field.Name, element.Kind())
		}
	}

	return nil
}

func encodeString(v reflect.Value, field reflect.StructField, buf *bytes.Buffer) error {
	switch v.Kind() {
	case reflect.String:
		s := v.String()
		if err := packet.WriteVarInt(buf, int32(len(s))); err != nil {
			return err
		}
		buf.WriteString(s)

	case reflect.Slice:
		s := v.Bytes()
		if err := packet.WriteVarInt(buf, int32(len(s))); err != nil {
			return err
		}
		buf.Write(s)

	default:
		return fmt.Errorf("field %s: value must be type string or byte slice, got %s", field.Name, v.Kind())
	}

	return nil
}
