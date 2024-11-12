package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"math/rand"

	"go.opentelemetry.io/otel"
)

var traceName string = "sample"
var (
	meter     = otel.Meter(traceName)
	helloType metric.Int64Counter
)

func init() {
	var err error
	helloType, err = meter.Int64Counter(
		"hello.message.type",
		metric.WithDescription("The number of hello message by roll value"),
		metric.WithUnit("{message}"))

	if err != nil {
		panic(err)
	}
}

func getHello(w http.ResponseWriter, r *http.Request) {
	var messageInit int
	ctx, span := otel.Tracer(traceName).Start(r.Context(), "getHello")
	defer span.End()
	w.WriteHeader(http.StatusOK)
	messageInit = rand.Intn(4) + 1
	switch messageInit {
	case 1:
		slog.ErrorContext(ctx, "Hello, World!")
	case 2:
		slog.WarnContext(ctx, "Hello, World!")
	case 3:
		slog.DebugContext(ctx, "Hello, World!")
	default:
		slog.InfoContext(ctx, "Hello, World!")
	}
	span.SetAttributes(attribute.Int("message.value", messageInit))
	helloType.Add(ctx, 1, metric.WithAttributes(
		attribute.Int("message.value", messageInit),
	))
	w.Write([]byte("Hello, World!"))
}

func getSleep(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(traceName).Start(r.Context(), "getSleep")
	defer span.End()
	sleepTime(ctx)
}

func sleepTime(ctx context.Context) {
	ctx, span := otel.Tracer(traceName).Start(ctx, "sleep")
	defer span.End()
	time.Sleep(time.Duration(rand.Intn(10)+1) * time.Microsecond * 10)
}
