package config

import (
	"log/slog"

	"github.com/m-mizutani/clog"
)

func LoggerConfig(level slog.Leveler) *slog.Logger {

	handler := handlerConfig(level)
	return slog.New(handler)
}

func handlerConfig(level slog.Leveler) *clog.Handler {
	handler := clog.New(
		clog.WithColor(true),
		clog.WithSource(true),
		clog.WithTimeFmt("2006-01-02T15:04:05.999Z07:00"),
		clog.WithLevel(level),
	)
	return handler
}
