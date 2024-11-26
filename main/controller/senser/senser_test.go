package senser

import (
	"context"
	"fmt"
	"log/slog"
	"senseregent/config"
	"testing"
	"time"
)

func TestSenserRun(t *testing.T) {
	slog.SetDefault(config.LoggerConfig(slog.LevelDebug))

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		time.Sleep(1 * time.Second)
		if err := Reset(ctx); err != nil {
			t.Fatalf("Reset Error")
		}
		time.Sleep(200 * time.Millisecond)
		if v, err := GetValue(ctx); err != nil {
			t.Fatalf("GetValue Error")
		} else {
			fmt.Printf("GetValue OK BME280 %v\n\n", v.BME280)
			fmt.Println(v.ToJson())
			fmt.Println(v.ToPromQL())
		}
		time.Sleep(2 * time.Second)
		if err := Stop(ctx); err != nil {
			t.Fatalf("Stop Error")
		}
		time.Sleep(500 * time.Millisecond)
		cancel()
	}()
	if err := Run(ctx); err != nil {
		t.Fatalf("Run Error")
	}
	time.Sleep(1 * time.Second)

}
