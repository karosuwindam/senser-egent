package i2csennser

import (
	"context"
	"errors"
	"log/slog"
	"senseregent/controller/sennser/i2c_sennser/bme280"
)

type I2CSennser struct {
	Value map[string]interface{}
}

type sennsertype struct {
	flag bool
	api  interface{}
}

type Bme280Value struct {
	Tmp   float64
	Hum   float64
	Press float64
}

var i2cSennserType map[string]sennsertype

func Init() error {
	slog.Debug("I2C Sennser Init")

	i2cSennserType = make(map[string]sennsertype)
	i2cSennserType["BME280"] = sennsertype{false, bme280.APIInit()}
	return nil
}

func Test(ctx context.Context) (err error) {
	slog.DebugContext(ctx, "i2c Test Start")

	err = nil
	//err にはエラーがあればエラーを登録する
	testErr := func(inErr error) {
		err = errors.Join(err, inErr)
	}

	for name, v := range i2cSennserType {
		switch name {
		case "BME280":
			api, ok := v.api.(*bme280.API)
			if !ok {
				i2cSennserType[name] = sennsertype{false, api}

				testErr(errors.New("BME280 API Type Error"))
				continue
			}
			slog.Info("BME280 Test Start")
			i2cSennserType[name] = sennsertype{api.Test(ctx), api}
		}
	}
	return
}

func SenserInit(ctx context.Context) (err error) {
	slog.DebugContext(ctx, "I2C SenserInit Start")

	err = nil
	//err にはエラーがあればエラーを登録する
	initErr := func(inErr error) {
		err = errors.Join(err, inErr)
	}
	for name, v := range i2cSennserType {
		switch name {
		case "BME280":
			if !v.flag {
				continue
			}
			api, ok := v.api.(*bme280.API)
			if !ok {
				i2cSennserType[name] = sennsertype{false, api}
				initErr(errors.New("BME280 API Type Error"))
				continue
			}
			if v.flag {
				slog.Info("BME280 Up Start")
				api.Up(ctx)
			}
		}
	}
	return
}

func SennserClose(ctx context.Context) (err error) {
	slog.DebugContext(ctx, "I2C SennserClose Start")

	err = nil
	//err にはエラーがあればエラーを登録する
	closeErr := func(inErr error) {
		err = errors.Join(err, inErr)
	}
	for name, v := range i2cSennserType {
		switch name {
		case "BME280":
			api, ok := v.api.(*bme280.API)
			if !ok {
				i2cSennserType[name] = sennsertype{false, api}
				closeErr(errors.New("BME280 API Type Error"))
				continue
			}
			if v.flag {
				slog.Info("BME280 Down Start")
				api.Down(ctx)
				i2cSennserType[name] = sennsertype{false, api}
			}
		}
	}
	return
}

func ReadValue(ctx context.Context) (value I2CSennser, err error) {
	slog.DebugContext(ctx, "I2C ReadValue Start")

	value = I2CSennser{
		Value: map[string]interface{}{},
	}
	readErr := func(inErr error) {
		err = errors.Join(err, inErr)
	}
	for name, v := range i2cSennserType {
		switch name {
		case "BME280":
			api, ok := v.api.(*bme280.API)
			if !ok {
				value.Value[name] = nil
				readErr(errors.New("BME280 API Type Error"))
				continue
			}
			if v.flag {
				slog.Info("BME280 Read Start")
				if err := api.ReadData(ctx); err != nil {
					value.Value[name] = nil
					readErr(err)
					continue
				}
				slog.Debug("BME280 Read Value", "Temp", api.Tmp, "Hum", api.Hum, "Press", api.Press)
				value.Value[name] = Bme280Value{
					Tmp:   api.Tmp,
					Hum:   api.Hum,
					Press: api.Press,
				}
			}
		}
	}
	return
}

func (value *I2CSennser) ReadBME280_value() Bme280Value {
	bme280Value, ok := value.Value["BME280"].(Bme280Value)
	if !ok {
		slog.Warn("BME280 Value Type Error")
		return Bme280Value{-1, -1, -1}
	}
	return bme280Value
}
