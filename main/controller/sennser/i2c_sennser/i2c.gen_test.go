package i2csennser

import (
	"context"
	"log/slog"
	"os"
	"testing"
)

func TestI2csennser(t *testing.T) {
	handler := slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelDebug},
	)
	logger := slog.New(handler)

	slog.SetDefault(logger)

	ctx := context.Background()
	if err := Init(); err != nil {
		t.Fatalf("Init Error")
	}
	if err := Test(ctx); err != nil {
		t.Fatalf("Test Error")
	}
	if err := SenserInit(ctx); err != nil {
		t.Fatalf("SenserInit Error")
	}
	v, err := ReadValue(ctx)
	if err != nil {
		t.Fatalf("ReadValue Error")
	}
	if a := v.ReadBME280_value(); a.Hum == -1 || a.Press == -1 || a.Tmp == -1 {
		t.Fatalf("ReadBME280_value Error")
	} else {
		t.Logf("ReadBME280_value OK %v", a)
	}
	if err := SennserClose(ctx); err != nil {
		t.Fatalf("SennserClose Error")
	}
}
