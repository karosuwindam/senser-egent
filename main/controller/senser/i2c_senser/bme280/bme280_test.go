package bme280

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"senseregent/controller/senser/i2c_senser/common"
	"testing"
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
