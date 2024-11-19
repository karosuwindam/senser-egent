package config

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/metric"

	sdklog "go.opentelemetry.io/otel/sdk/log"

	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func initConn(url string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}
	return conn, err
}

func initTracerProvider(ctx context.Context, res *resource.Resource, v interface{}) (*sdktrace.TracerProvider, error) {
	var traceExporter *otlptrace.Exporter
	var err error
	if conn, ok := v.(*grpc.ClientConn); ok {
		traceExporter, err = otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
		if err != nil {
			return nil, fmt.Errorf("failed to create trace exporter: %w", err)
		}
	} else if url, ok := v.(string); ok {
		traceExporter, err = otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint(url),
		)
	}
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return tracerProvider, nil
}

func initMeterProvider(ctx context.Context, res *resource.Resource, v interface{}) (*metric.MeterProvider, error) {
	var meterProvider *metric.MeterProvider
	if conn, ok := v.(*grpc.ClientConn); ok {
		metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
		if err != nil {
			return nil, fmt.Errorf("failed to create metric exporter: %w", err)
		}
		meterProvider = metric.NewMeterProvider(
			metric.WithResource(res),
			metric.WithReader(metric.NewPeriodicReader(metricExporter,
				metric.WithInterval(3*time.Second))),
		)
	} else if url, ok := v.(string); ok {
		metricExporter, err := otlpmetrichttp.New(ctx,
			otlpmetrichttp.WithEndpoint(url),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create metric exporter: %w", err)
		}
		meterProvider = metric.NewMeterProvider(
			metric.WithResource(res),
			metric.WithReader(metric.NewPeriodicReader(metricExporter,
				metric.WithInterval(3*time.Second))),
		)
	}

	return meterProvider, nil
}

func initLoggerProvider(ctx context.Context, res *resource.Resource, v interface{}) (*sdklog.LoggerProvider, error) {
	var logExporter sdklog.Exporter
	var err error
	if conn, ok := v.(*grpc.ClientConn); ok {
		logExporter, err = otlploggrpc.New(ctx,
			otlploggrpc.WithGRPCConn(conn),
		)
	} else if url, ok := v.(string); ok {
		logExporter, err = otlploghttp.New(ctx,
			otlploghttp.WithEndpointURL(url),
		)
	}
	if err != nil {
		return nil, err
	}

	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
	)
	return loggerProvider, nil
}

var shutdownFuncs []func(context.Context) error

func TracerStart(urldata, serviceName string, ctx context.Context) error {

	if !TraData.TracerUse {
		return nil
	}
	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return err
	}
	if TraData.GrpcOn {
		conn, err := initConn(urldata)
		if err != nil {
			return err
		}
		if shutdownTracer, err := initTracerProvider(ctx, res, conn); err != nil {
			return err
		} else {
			shutdownFuncs = append(shutdownFuncs, shutdownTracer.Shutdown)
		}
		if shutdownMeter, err := initMeterProvider(ctx, res, conn); err != nil {
			return err
		} else {
			shutdownFuncs = append(shutdownFuncs, shutdownMeter.Shutdown)
		}
		if shutdownlogger, err := initLoggerProvider(ctx, res, "http://"+TraData.HttpURL+"/v1/logs"); err != nil {
			return err
		} else {
			global.SetLoggerProvider(shutdownlogger)
			shutdownFuncs = append(shutdownFuncs, shutdownlogger.Shutdown)
		}
	} else {
		if shutdownTracer, err := initTracerProvider(ctx, res, TraData.HttpURL); err != nil {
			return err
		} else {
			shutdownFuncs = append(shutdownFuncs, shutdownTracer.Shutdown)
		}
		if shutdownMeter, err := initMeterProvider(ctx, res, TraData.HttpURL); err != nil {
			return err
		} else {
			shutdownFuncs = append(shutdownFuncs, shutdownMeter.Shutdown)
		}
		if shutdownlogger, err := initLoggerProvider(ctx, res, "http://"+TraData.HttpURL+"/v1/logs"); err != nil {
			return err
		} else {
			global.SetLoggerProvider(shutdownlogger)
			shutdownFuncs = append(shutdownFuncs, shutdownlogger.Shutdown)
		}
	}
	logger := slog.New(
		slogmulti.Fanout(
			handlerConfig(slog.LevelInfo),
			otelslog.NewHandler(serviceName),
		),
	)
	slog.SetDefault(logger)

	// l := otelslog.NewLogger(serviceName)
	// l.InfoContext(ctx, "gageg")
	// slog.SetDefault(otelslog.NewLogger(serviceName))
	// logeger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// logeger.With()

	return nil
}

func TracerStop(ctx context.Context) error {
	if !TraData.TracerUse {
		return nil
	}
	var err_ch chan error = make(chan error, 1)
	go func(ctx context.Context) {
		var err error
		for _, f := range shutdownFuncs {
			err = errors.Join(err, f(ctx))
		}
		shutdownFuncs = nil
		err_ch <- err
	}(ctx)
	select {
	case err := <-err_ch:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(5 * time.Second):
		return errors.New("timeout")
	}
}

type TraceSet struct {
	Otel trace.Span
}

func TracerS(ctx context.Context,
	processName, spanName string,
	opts ...trace.SpanStartOption) (context.Context, *TraceSet) {
	ctx, ts := otel.Tracer(processName).Start(ctx, spanName, opts...)
	return ctx, &TraceSet{ts}
}

func (t *TraceSet) End() {
	if t.Otel != nil {
		t.Otel.End()
	}
}

func (t *TraceSet) SetAttributes(kv ...attribute.KeyValue) {
	if t.Otel != nil {
		t.Otel.SetAttributes(kv...)
	}
}
