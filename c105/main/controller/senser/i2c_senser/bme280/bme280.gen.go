package bme280

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"senseregent/controller/senser/i2c_senser/common"
	"sync"
	"time"
)

var (
	shutdown chan bool
	oneshut  chan bool
	done     chan bool
	reset    chan bool
)

type contextData struct {
	value Bme280_Vaule
	flag  bool
	mu    sync.Mutex
}

var datastore contextData

func Init(ctx context.Context) error {
	var err error = nil
	i2c = common.Init(BME280, I2C_BUS)

	shutdown = make(chan bool, 1)
	oneshut = make(chan bool, 1)
	done = make(chan bool, 1)
	reset = make(chan bool, 1)

	// senser Init
	up(ctx)
	defer down(ctx)
	for i := 0; i < 3; i++ {
		if err = Test(ctx); err == nil {
			slog.InfoContext(ctx, "BME280 start check ok")
			break
		} else {
			slog.ErrorContext(ctx, "BME280 start test error", "count", i+1, "err", err)
		}
	}
	if err != nil {
		return errors.New("not Init Error for AM2320")
	}
	datastore.ON()
	slog.InfoContext(ctx, "BME280 Init OK")
	return nil
}

func Run(ctx context.Context) error {
	var calib bme280_cal

	slog.InfoContext(ctx, "BME280 loop start")
	oneshut <- true
	up(ctx)
loop:
	for {
		select {
		case <-ctx.Done():
			slog.ErrorContext(ctx, "context Done")
			break loop
		case <-reset:
			datastore.OFF()
			down(ctx)
			up(ctx)
			for i := 0; i < 3; i++ {
				if err := Test(ctx); err == nil {
					calib = calibRead(ctx)
					datastore.ON()
					break
				}
				time.Sleep(100 * time.Millisecond)
			}
		case <-shutdown:
			slog.InfoContext(ctx, "BME280 run loop shutdown")
			done <- true
			break loop
		case <-oneshut:
			if !datastore.Check() {
				break
			}
			calib = calibRead(ctx)
			press, tmp, hum := readSenserData(ctx, calib)
			datastore.SetValue(Bme280_Vaule{press, tmp, hum})
			fmt.Printf("press=%v,tmp=%v,hum=%v\n", press, tmp, hum)
		case <-time.After(100 * time.Millisecond):
			if !datastore.Check() {
				break
			}
			press, tmp, hum := readSenserData(ctx, calib)
			datastore.SetValue(Bme280_Vaule{press, tmp, hum})
			fmt.Printf("press=%v,tmp=%v,hum=%v\n", press, tmp, hum)
		}
	}
	down(ctx)
	slog.InfoContext(ctx, "BME280 loop stop")
	return nil
}

func Stop(ctx context.Context) error {
	shutdown <- true
	select {
	case <-ctx.Done():
		return errors.New("context Done")
	case <-done:
		break
	case <-time.After(1 * time.Second):
		// msg := "shutdown time out"
		return errors.New("shutdown time out")
	}
	return nil
}
func Test(ctx context.Context) error {
	errch := make(chan error, 1)
	var err error = nil
	go func(ctx context.Context) {
		if readICIDCheck(ctx) {
			errch <- errors.New("BME280 Test check error")
		} else {
			errch <- nil
		}
	}(ctx)

	select {
	case <-ctx.Done():
		err = errors.New("context Done")
	case <-time.After(5 * time.Second):
		err = fmt.Errorf("5second over time")
	case err = <-errch:
		break
	}
	return err
}

func Reset() {
	if len(reset) > 0 {
		return
	}
	reset <- true
}

func ReadValue() (Bme280_Vaule, bool) {
	v := datastore.Value()
	f := datastore.Check()
	return v, f
}

func (c *contextData) ON() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.flag = true
}

func (c *contextData) OFF() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.flag = false
}

func (c *contextData) Check() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.flag
}

func (c *contextData) SetValue(v Bme280_Vaule) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value = v
}

func (c *contextData) Value() Bme280_Vaule {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}
