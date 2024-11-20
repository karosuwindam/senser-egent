package controller

import (
	"context"
	"fmt"
	"log/slog"
	"senseregent/controller/sennser"
	"time"
)

var shutdown chan struct{}
var done chan struct{}

func Init() error {

	shutdown = make(chan struct{}, 1)
	done = make(chan struct{}, 1)
	return nil
}

func Run(ctx context.Context) error {
	slog.InfoContext(ctx, "Controller Run Start")
	go func(ctx context.Context) {
		if err := sennser.Run(ctx); err != nil {
			slog.ErrorContext(ctx, "Sennser Run error", "error", err)
		}
	}(ctx)
	select {
	case <-shutdown:
		if err := sennser.Stop(ctx); err != nil {
			slog.ErrorContext(ctx, "Sennser Stop error", "error", err)
		}
		done <- struct{}{}
	case <-ctx.Done():
		slog.ErrorContext(ctx, "Sennser Run Stop by context")
	}
	slog.InfoContext(ctx, "Controller Run Stop")
	return nil
}

func Stop(ctx context.Context) error {
	slog.InfoContext(ctx, "Controller Stop Start")
	shutdown <- struct{}{}
	select {
	case <-done:
		break
	case <-time.After(5 * time.Second):
		return fmt.Errorf("time over 5 sec")
	}
	return nil
}
