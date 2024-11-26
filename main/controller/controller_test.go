package controller

import (
	"context"
	"fmt"
	"log/slog"
	"senseregent/config"
	"testing"
	"time"
)

func TestController(t *testing.T) {
	slog.SetDefault(config.LoggerConfig(slog.LevelInfo))
	if err := Init(); err != nil {
		t.Errorf("Init error %v", err)
	}
	ctx := context.Background()
	go func() {
		api := NewAPI()
		for i := 0; i < 3; i++ {
			if v, err := api.ReadValue(ctx); err != nil {
				t.Errorf("ReadValue error %v", err)
			} else {
				fmt.Println("BME280", v.BME280)
				fmt.Println(v.ToValueType())
				fmt.Println(v.ToJson())
			}
			time.Sleep(300 * time.Millisecond)
		}
		if err := api.ResetSenser(ctx); err != nil {
			t.Errorf("ResetSenser error %v", err)
		}
		for i := 0; i < 3; i++ {
			if v, err := api.ReadValue(ctx); err != nil {
				t.Errorf("ReadValue error %v", err)
			} else {
				fmt.Println("BME280", v.BME280)
				fmt.Println(v.ToJson())
			}
			time.Sleep(300 * time.Millisecond)
		}
		time.Sleep(4 * time.Second)
		Stop(ctx)
	}()
	if err := Run(ctx); err != nil {
		t.Errorf("Run error %v", err)
	}
}
