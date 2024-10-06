package common

import (
	"context"
	"sync"
	"time"

	"github.com/davecheney/i2c"
	"github.com/pkg/errors"
)

type I2CBUS struct {
	addr uint8
	bus  int
}

var i2cMu sync.Mutex

func Init(addr uint8, bus int) *I2CBUS {
	return &I2CBUS{addr: addr, bus: bus}
}

func (t *I2CBUS) WriteByte(ctx context.Context) error {
	i2cMu.Lock()
	defer i2cMu.Unlock()

	done := make(chan struct{}, 1)
	errch := make(chan error, 1)
	i2cd, ok := readI2cContext(ctx)
	if !ok {
		return errors.New("error context data not fond I2C Data")
	}
	go func(command, data byte) {
		defer func() { done <- struct{}{} }()
		if err := t.writeByte(command, data); err != nil {
			errch <- err
		}
	}(i2cd.Command, i2cd.Data)
	select {
	case <-ctx.Done():
		return errors.New("context Done")
	case <-time.After(time.Second * 5):
		return errors.New("I2c Write timeout")
	case err := <-errch:
		return err
	case <-done:
		break
	}
	return nil
}

func (t *I2CBUS) ReadByte(ctx context.Context) ([]byte, error) {

	i2cMu.Lock()
	defer i2cMu.Unlock()

	done := make(chan struct{}, 1)
	errch := make(chan error, 1)
	i2cd, ok := readI2cContext(ctx)
	if !ok {
		return []byte{}, errors.New("error context data not fond I2C Data")
	}

	buf := make([]byte, i2cd.ReadSize)
	go func(command byte, size int) {
		defer func() { done <- struct{}{} }()
		tmp, err := t.readByte(command, size)
		buf = tmp
		if err != nil {
			errch <- err
		}
	}(i2cd.Command, i2cd.ReadSize)

	select {
	case <-ctx.Done():
		return buf, errors.New("context Done")
	case <-time.After(time.Second * 5):
		return buf, errors.New("I2c Read timeout")
	case err := <-errch:
		return buf, err
	case <-done:
		break
	}
	return buf, nil
}

func (t *I2CBUS) writeByte(command, data byte) error {
	i2c, err := i2c.New(t.addr, t.bus)
	if err != nil {
		return errors.Wrapf(err, "i2c.New(%v,%v)", t.addr, t.bus)
	}
	defer i2c.Close()
	_, err = i2c.Write([]byte{command, data})
	if err != nil {
		return errors.Wrapf(err, "i2c.Write(%v,%v)", command, data)
	}
	return nil
}

func (t *I2CBUS) readByte(command byte, size int) ([]byte, error) {
	buf := make([]byte, size)
	i2c, err := i2c.New(t.addr, t.bus)
	if err != nil {
		return buf, errors.Wrapf(err, "i2c.New(%v,%v)", t.addr, t.bus)
	}
	defer i2c.Close()
	_, err = i2c.Write([]byte{command})
	if err != nil {
		return buf, errors.Wrapf(err, "i2c.Write(%v)", command)
	}
	_, err = i2c.Read(buf)
	if err != nil {
		return buf, errors.Wrap(err, "i2c.Read()")
	}
	return buf, nil

}
