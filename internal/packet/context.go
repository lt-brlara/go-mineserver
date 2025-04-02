package packet

import (
	"context"
)

type key int

const (
	lengthKey key = iota + 1
	packetIdKey
)

func NewContext(ctx context.Context, length int32, id byte) context.Context {
	ctx = context.WithValue(ctx, lengthKey, length)
	ctx = context.WithValue(ctx, packetIdKey, id)
	return ctx
}

func IdFromContext(ctx context.Context) (byte, bool) {
	id, ok := ctx.Value(packetIdKey).(byte)
	return id, ok
}

func LengthFromContext(ctx context.Context) (int32, bool) {
	length, ok := ctx.Value(lengthKey).(int32)
	return length, ok
}
