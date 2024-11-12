package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var grpcflag bool = true

// var otlpUrl string = "localhost:4317"

var otlpUrl string = "otel.bookserver.home:4317"

func initConn(url string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func initTracerProvider(ctx context.Context, res *resource.Resource, opt interface{}) (
	*sdktrace.TracerProvider, error) {
	var traceExporter *otlptrace.Exporter
	var err error
	if conn, ok := opt.(*grpc.ClientConn); ok {
		traceExporter, err = otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	} else if url, ok := opt.(string); ok {
		traceExporter, err = otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint(url),
		)
	}
	if err != nil {
		return nil, err
	}
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter,
		sdktrace.WithBatchTimeout(time.Second),
	)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return tracerProvider, nil
}

func initMeterProvider(ctx context.Context, res *resource.Resource, opt interface{}) (
	*sdkmetric.MeterProvider, error) {
	var meterProvider *sdkmetric.MeterProvider
	var meterExporter sdkmetric.Exporter
	var err error

	if conn, ok := opt.(*grpc.ClientConn); ok {
		meterExporter, err = otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	} else if url, ok := opt.(string); ok {
		meterExporter, err = otlpmetrichttp.New(ctx,
			otlpmetrichttp.WithEndpoint(url),
		)
	}
	if err != nil {
		return nil, err
	}
	meterProvider = sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(meterExporter,
			sdkmetric.WithInterval(3*time.Second))))

	return meterProvider, nil
}

func initLoggerProvider(ctx context.Context, res *resource.Resource, opt interface{}) (
	*sdklog.LoggerProvider, error) {
	var logerExporter sdklog.Exporter
	var err error
	if conn, ok := opt.(*grpc.ClientConn); ok {
		logerExporter, err = otlploggrpc.New(ctx,
			otlploggrpc.WithGRPCConn(conn),
		)
	} else if url, ok := opt.(string); ok {
		logerExporter, err = otlploghttp.New(ctx,
			otlploghttp.WithEndpointURL(url),
		)
	}
	if err != nil {
		return nil, err
	}

	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logerExporter)),
	)
	return loggerProvider, nil
}

func initResourc(ctx context.Context, servicName string) (*resource.Resource, error) {
	// return resource.New(ctx,
	// 	resource.WithAttributes(
	// 		semconv.ServiceNameKey.String(servicName),
	// 		semconv.ServiceVersion("0.1.0"),
	// 	),
	// )
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(servicName),
			semconv.ServiceVersion("0.1.0"),
		))
}

func setupOTel(ctx context.Context, serviceName string) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}
	res, err := initResourc(ctx, serviceName)
	if err != nil {
		handleErr(err)
		return
	}

	var tracerProvider *sdktrace.TracerProvider
	var meterProvider *sdkmetric.MeterProvider
	var loggerProvider *sdklog.LoggerProvider
	var errtmp error
	if grpcflag {
		conn, errtmp := initConn(otlpUrl)
		if errtmp != nil {
			handleErr(errtmp)
			return
		}
		tracerProvider, errtmp = initTracerProvider(ctx, res, conn)
		if errtmp != nil {
			handleErr(errtmp)
			return
		}
		meterProvider, errtmp = initMeterProvider(ctx, res, conn)
		if errtmp != nil {
			handleErr(errtmp)
			return
		}
		loggerProvider, errtmp = initLoggerProvider(ctx, res, conn)
		if errtmp != nil {
			handleErr(errtmp)
			return
		}
	} else {
		tracerProvider, errtmp = initTracerProvider(ctx, res, otlpUrl)
		if errtmp != nil {
			handleErr(errtmp)
			return
		}
		meterProvider, errtmp = initMeterProvider(ctx, res, otlpUrl)
		if errtmp != nil {
			handleErr(errtmp)
			return
		}
		loggerProvider, errtmp = initLoggerProvider(ctx, res, otlpUrl)
		if errtmp != nil {
			handleErr(errtmp)
			return
		}

	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)
	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)
	logger := slog.New(
		slogmulti.Fanout(
			slog.NewTextHandler(os.Stdout, nil),
			otelslog.NewHandler(serviceName),
		),
	)
	slog.SetDefault(logger)

	return
}
