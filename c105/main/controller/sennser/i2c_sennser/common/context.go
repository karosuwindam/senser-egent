package common

import "context"

type I2c struct {
	Command  byte
	Data     byte
	ReadSize int
}

func WriteI2cContext(ctx context.Context, i I2c) context.Context {
	return context.WithValue(ctx, I2c{}, i)
}

func readI2cContext(ctx context.Context) (I2c, bool) {
	i, ok := ctx.Value(I2c{}).(I2c)
	return i, ok
}
