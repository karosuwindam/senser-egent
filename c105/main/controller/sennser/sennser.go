package sennser

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"senseregent/controller/sennser/common"
	i2csennser "senseregent/controller/sennser/i2c_sennser"
	"sync"
	"time"
)

var shutdown chan struct{}
var done chan struct{}
var reset chan struct{}

func init() {
	if err := i2csennser.Init(); err != nil {
		panic(err)
	}
	shutdown = make(chan struct{}, 1)
	done = make(chan struct{}, 1)
	reset = make(chan struct{}, 1)
}

func Run(ctx context.Context) error {
	slog.InfoContext(ctx, "Sennser Run Start")
	teststop := make(chan struct{}, 1)
	go func(ctx context.Context) {
		slog.DebugContext(ctx, "Sennser Test Check Start")
		var oneshut chan struct{} = make(chan struct{}, 1)
		oneshut <- struct{}{}
		for {
			select {
			case <-oneshut:
				if err := testSennser(ctx); err != nil {
					slog.ErrorContext(ctx, "testSenser error", "error", err)
				}
			case <-reset:
				if err := closeSennser(ctx); err != nil {
					slog.ErrorContext(ctx, "closeSennser error", "error", err)
				}
				oneshut <- struct{}{}
			case <-teststop:
				slog.DebugContext(ctx, "Sennser Test Check Stop")
				return
			case <-ctx.Done():
				slog.ErrorContext(ctx, "Sennser Test Check Stop")
				return
			}
		}
	}(ctx)
loop:
	for {
		select {
		case <-shutdown:
			slog.DebugContext(ctx, "Sennser Run Shutdown Start")
			if err := closeSennser(ctx); err != nil {
				slog.ErrorContext(ctx, "closeSennser error", "error", err)
			}
			teststop <- struct{}{}
			done <- struct{}{}
			break loop
		case <-ctx.Done():
			break loop
		case <-time.After(200 * time.Millisecond):
			v, err := readSennser(ctx)
			if err != nil {
				slog.ErrorContext(ctx, "readSennser error", "error", err)
				continue
			}
			if err := setValue(ctx, v); err != nil {
				slog.ErrorContext(ctx, "setValue error", "error", err)
			}
			slog.DebugContext(ctx, "readSennser", "value", v)
		}
	}
	slog.InfoContext(ctx, "Sennser Run Stop")
	return nil
}

func Stop(ctx context.Context) error {
	slog.InfoContext(ctx, "Sennser Stop Start")
	shutdown <- struct{}{}
	select {
	case <-done:
	case <-ctx.Done():
		return errors.New("Stop Timeout")
	case <-time.After(5 * time.Second):
		return errors.New("Stop Timeout over 5s")
	}
	return nil
}

func Reset(ctx context.Context) error {
	if len(reset) > 0 {
		return errors.New("Reset Already")
	}
	slog.InfoContext(ctx, "Sennser Reset")
	reset <- struct{}{}
	return nil
}

type SennserValue struct {
	BME280 *i2csennser.Bme280Value
}

type ValueType struct {
	Senser      string `json:"Senser"`
	Type        string `json:"Type"`
	Data        string `json:"Data"`
	help        string
	types       string
	promeQLName string
}

func GetValue(ctx context.Context) (SennserValue, error) {
	var output SennserValue = SennserValue{}
	v, err := getValue(ctx)
	if err != nil {
		return output, err
	}
	if len(v) == 0 {
		return output, errors.New("Sennser Value is Empty")
	}
	for name, val := range v {
		switch name {
		case "BME280":
			bme280, ok := val.(i2csennser.Bme280Value)
			if !ok {
				return SennserValue{}, errors.New("BME280 Value Type Error")
			}
			output.BME280 = &bme280
		}
	}

	return output, nil
}

func (v *SennserValue) ToValueType() []ValueType {
	var output []ValueType
	if v.BME280 != nil {
		output = append(output, v.toBME280Type()...)
	}
	return output
}

func (v *SennserValue) ToJson() string {
	tmpValue := v.ToValueType()
	json, err := json.Marshal(tmpValue)
	if err != nil {
		return ""
	}
	return string(json)
}

func (v *SennserValue) toBME280Type() []ValueType {
	var output []ValueType
	tmp := fmt.Sprintf("%.2f", v.BME280.Tmp)
	output = append(output, ValueType{
		Senser:      "BME280",
		Type:        "tmp",
		Data:        tmp,
		help:        common.PROMQL_HELP_TMP,
		types:       common.PROMQLTYPE_GAUGE,
		promeQLName: common.PROMQLNAME_BME280,
	})
	hum := fmt.Sprintf("%.2f", v.BME280.Hum)
	output = append(output, ValueType{
		Senser:      "BME280",
		Type:        "hum",
		Data:        hum,
		help:        common.PROMQL_HELP_HUM,
		types:       common.PROMQLTYPE_GAUGE,
		promeQLName: common.PROMQLNAME_BME280,
	})
	press := fmt.Sprintf("%.2f", v.BME280.Press)
	output = append(output, ValueType{
		Senser:      "BME280",
		Type:        "press",
		Data:        press,
		help:        common.PROMQL_HELP_PRESS,
		types:       common.PROMQLTYPE_GAUGE,
		promeQLName: common.PROMQLNAME_BME280,
	})
	return output
}

func (v *SennserValue) ToPromQL() string {
	var output string
	if v.BME280 != nil {
		if output == "" {
			output = v.toBME280PromQL()
		} else {
			output += v.toBME280PromQL()
		}
	}
	return output
}

func (v *SennserValue) toBME280PromQL() string {
	var output string
	if v.BME280 != nil {
		velue := v.toBME280Type()
		for _, v := range velue {
			output += v.promqlHelp()
			output += v.promqlType()
			output += v.promqlValue("BME280")
		}
	}
	return output
}

func (v *ValueType) promqlHelp() string {
	return fmt.Sprintf("# HELP %s %s\n", v.promeQLName, v.help)
}

func (v *ValueType) promqlType() string {
	return fmt.Sprintf("# TYPE %s %s\n", v.promeQLName, v.types)
}

func (v *ValueType) promqlValue(sennserName string) string {
	return fmt.Sprintf("%s{type=\"%s\", sennser=\"%s\"} %s\n",
		"senser"+"_"+v.Type+"_value", v.Type, sennserName, v.Data,
	)
}

var tmpvalue = map[string]interface{}{}
var tmpSync sync.Mutex

func setValue(ctx context.Context, value map[string]interface{}) error {
	slog.DebugContext(ctx, "setValue", "value", value)
	tmpSync.Lock()
	defer tmpSync.Unlock()
	tmpvalue = value
	return nil
}

func getValue(ctx context.Context) (map[string]interface{}, error) {
	slog.DebugContext(ctx, "GetValue", "value", tmpvalue)
	tmpSync.Lock()
	defer tmpSync.Unlock()
	return tmpvalue, nil
}

func testSennser(ctx context.Context) (err error) {
	slog.DebugContext(ctx, "Test Senser Start")

	err = nil
	testErr := func(inErr error) {
		err = errors.Join(err, inErr)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		if err := i2csennser.Test(ctx); err != nil {
			testErr(err)
		} else {
			i2csennser.SenserInit(ctx)
		}
	}(ctx)
	wg.Wait()
	return
}

func readSennser(ctx context.Context) (value map[string]interface{}, err error) {
	slog.DebugContext(ctx, "Read Senser Start")

	err = nil
	value = map[string]interface{}{}
	readErr := func(inErr error) {
		err = errors.Join(err, inErr)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		slog.DebugContext(ctx, "readSennser I2C ReadValue Start")
		v, err := i2csennser.ReadValue(ctx)
		if err != nil {
			readErr(err)
			return
		}
		for name, val := range v.Value {
			value[name] = val
		}
	}(ctx)
	wg.Wait()
	return
}

func closeSennser(ctx context.Context) (err error) {
	slog.DebugContext(ctx, "Close Senser Start")

	err = nil
	closeErr := func(inErr error) {
		err = errors.Join(err, inErr)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		if err := i2csennser.SennserClose(ctx); err != nil {
			closeErr(err)
		}
	}(ctx)
	wg.Wait()
	return
}
