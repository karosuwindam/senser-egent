package i2csenser

import (
	"context"
	"errors"
	"log/slog"
	"senseregent/config"
	"senseregent/controller/senser/i2c_senser/bme280"
)

type I2CSenser struct {
	Value map[string]interface{}
}

type sensertype struct {
	flag bool
	api  interface{}
}

type Bme280Value struct {
	Tmp   float64
	Hum   float64
	Press float64
}

var i2cSenserType map[string]sensertype

func Init() error {
	slog.Debug("I2C Senser Init")

	i2cSenserType = make(map[string]sensertype)
	i2cSenserType["BME280"] = sensertype{false, bme280.APIInit()}
	return nil
}

func Test(ctx context.Context) (err error) {
	ctx, span := config.TracerS(ctx, "Test", "i2c Test")
	defer span.End()
	slog.DebugContext(ctx, "i2c Test Start")

	err = nil
	//err にはエラーがあればエラーを登録する
	testErr := func(inErr error) {
		err = errors.Join(err, inErr)
	}

	for name, v := range i2cSenserType {
		switch name {
		case "BME280":
			api, ok := v.api.(*bme280.API)
			if !ok {
				i2cSenserType[name] = sensertype{false, api}

				testErr(I2C_Type_Error_BME280)
				continue
			}
			slog.Info("I2C BME280 Test Start")
			i2cSenserType[name] = sensertype{api.Test(ctx), api}
		}
	}
	return
}

func SenserInit(ctx context.Context) (err error) {
	ctx, span := config.TracerS(ctx, "SenserInit", "i2c SenserInit")
	defer span.End()
	slog.DebugContext(ctx, "I2C SenserInit Start")

	err = nil
	//err にはエラーがあればエラーを登録する
	initErr := func(inErr error) {
		err = errors.Join(err, inErr)
	}
	for name, v := range i2cSenserType {
		switch name {
		case "BME280":
			if !v.flag {
				continue
			}
			api, ok := v.api.(*bme280.API)
			if !ok {
				i2cSenserType[name] = sensertype{false, api}
				initErr(I2C_Type_Error_BME280)
				continue
			}
			slog.Info("I2C BME280 Up Start")
			api.Up(ctx)
		}
	}
	return
}

func SenserClose(ctx context.Context) (err error) {
	ctx, span := config.TracerS(ctx, "SenserClose", "i2c SenserClose")
	defer span.End()
	slog.DebugContext(ctx, "I2C SenserClose Start")

	err = nil
	//err にはエラーがあればエラーを登録する
	closeErr := func(inErr error) {
		err = errors.Join(err, inErr)
	}
	for name, v := range i2cSenserType {
		switch name {
		case "BME280":
			if !v.flag {
				continue
			}
			api, ok := v.api.(*bme280.API)
			if !ok {
				i2cSenserType[name] = sensertype{false, api}
				closeErr(I2C_Type_Error_BME280)
				continue
			}
			slog.Info("I2C BME280 Down Start")
			api.Down(ctx)
			i2cSenserType[name] = sensertype{false, api}
		}
	}
	return
}

func ReadValue(ctx context.Context) (value I2CSenser, err error) {
	ctx, span := config.TracerS(ctx, "ReadValue", "i2c ReadValue")
	defer span.End()
	slog.DebugContext(ctx, "I2C ReadValue Start")

	value = I2CSenser{
		Value: map[string]interface{}{},
	}
	readErr := func(inErr error) {
		err = errors.Join(err, inErr)
	}
	for name, v := range i2cSenserType {
		switch name {
		case "BME280":
			if !v.flag {
				continue
			}
			api, ok := v.api.(*bme280.API)
			if !ok {
				value.Value[name] = nil
				readErr(I2C_Type_Error_BME280)
				continue
			}
			slog.Info("I2C BME280 Read Start")
			if err := api.ReadData(ctx); err != nil {
				value.Value[name] = nil
				readErr(err)
				continue
			}
			slog.Debug("I2C BME280 Read Value", "senser", bme280.SENSER_NAME, "Temp", api.Tmp, "Hum", api.Hum, "Press", api.Press)
			value.Value[name] = Bme280Value{
				Tmp:   api.Tmp,
				Hum:   api.Hum,
				Press: api.Press,
			}
		}
	}
	return
}

func (value *I2CSenser) ReadBME280_value() Bme280Value {
	bme280Value, ok := value.Value["BME280"].(Bme280Value)
	if !ok {
		slog.Warn("I2C BME280 Value Type Error")
		return Bme280Value{-1, -1, -1}
	}
	return bme280Value
}