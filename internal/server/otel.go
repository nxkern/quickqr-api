package server

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

// SetupOTelSDK initializes OpenTelemetry providers for tracing, metrics, and logging using OTLP HTTP exporters.
// It returns a cleanup function to gracefully shut down the providers.
func SetupOTelSDK(ctx context.Context) (func(context.Context) error, error) {
	var shutdownFuncs []func(context.Context) error
	var err error

	shutdown := func(ctx context.Context) error {
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

	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	tracerProvider, err := newTracerProvider(ctx)
	if err != nil {
		handleErr(err)
		return shutdown, err
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	meterProvider, err := newMeterProvider(ctx)
	if err != nil {
		handleErr(err)
		return shutdown, err
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	loggerProvider, err := newLoggerProvider(ctx)
	if err != nil {
		handleErr(err)
		return shutdown, err
	}
	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	otel.SetLoggerProvider(loggerProvider)

	return shutdown, err

}

// newPropagator creates a composite propagator for context propagation.
func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

// newTracerProvider initializes OTLP HTTP trace provider.
func newTracerProvider(ctx context.Context) (*trace.TracerProvider, error) {
	traceExporter, err := otlptracehttp.New(ctx)
	if err != nil {
		return nil, err
	}
	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter, trace.WithBatchTimeout(5*time.Second)),
	)
	return tracerProvider, nil
}

// newMeterProvider initializes OTLP HTTP metric provider.
func newMeterProvider(ctx context.Context) (*metric.MeterProvider, error) {
	metricsExporter, err := otlpmetrichttp.New(ctx)
	if err != nil {
		return nil, err
	}
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricsExporter, metric.WithInterval(1*time.Minute))),
	)
	return meterProvider, nil
}

// newLoggerProvider initializes OTLP HTTP logger provider.
func newLoggerProvider(ctx context.Context) (*log.LoggerProvider, error) {
	logExporter, err := otlploghttp.New(ctx)

	if err != nil {
		return nil, err
	}

	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
	)
	return loggerProvider, nil
}
