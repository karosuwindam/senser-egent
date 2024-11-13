package bme280

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"senseregent/controller/sennser/i2c_sennser/common"
	"testing"
	"time"
)

func TestBme280Read(t *testing.T) {
	handler := slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelDebug},
	)
	logger := slog.New(handler)

	slog.SetDefault(logger)

	i2c = common.Init(BME280, I2C_BUS)
	ctx := context.Background()

	if readICIDCheck(ctx) {
		if readStatus(ctx) != CtrMeasReg_Sleep {
			t.Fatal("senser not Sleep")
		}
		up(ctx)

		if readStatus(ctx) != CtrMeasReg_Normal {
			t.Fatal("senser not Normal")
		}
		defer down(ctx)
		cal := calibRead(ctx)
		press, temp, hum := readSenserData(ctx, cal)
		fmt.Println("temp", temp, "hum", hum, "press", press)
	} else {
		t.Fatal("senser read error")
	}
}

func TestBme280Runloop(t *testing.T) {
	handler := slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelInfo},
	)
	logger := slog.New(handler)

	slog.SetDefault(logger)
	api := APIInit()
	ctx := context.Background()
	api.Up(ctx)
	api.CalibRead(ctx)
	go func(ctx context.Context) {
		for {
			if err := api.ReadData(ctx); err != nil {
				fmt.Println("loopstop")
				return
			}
			fmt.Println("temp", api.Tmp, "hum", api.Hum, "press", api.Press)
			time.Sleep(200 * time.Millisecond)
		}
	}(ctx)
	time.Sleep(3 * time.Second)
	api.Down(ctx)
	time.Sleep(1 * time.Second)
	if len(api.calib.hum) != 0 {
		t.Fatalf("Not clear calib by down ")
	}
}
