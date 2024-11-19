package bme280

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestBme280Runloop(t *testing.T) {
	handler := slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelInfo},
	)
	logger := slog.New(handler)

	slog.SetDefault(logger)
	api := APIInit()
	ctx := context.Background()
	if api.Test(ctx) == false {
		t.Fatalf("BME280 Test Error")
	}
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
