package instrumentation

import (
	"context"
	"io"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
)

func NewTracer(service string) io.Closer {
	ctx := context.Background()
	exp, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpoint("localhost:4318"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		panic(err)
	}

	resource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(service),
		),
	)
	if err != nil {
		panic(err)
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(resource),
		trace.WithSampler(trace.AlwaysSample()),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return &tracer{
		tracerProvider: tracerProvider,
	}
}

type tracer struct {
	tracerProvider *trace.TracerProvider
}

func (t *tracer) Close() error {
	return t.tracerProvider.Shutdown(context.Background())
}

var GRPC = []grpc.ServerOption{
	grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor(
		otelgrpc.WithPropagators(
			propagation.TraceContext{},
		),
	)),
	grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor(
		otelgrpc.WithPropagators(
			propagation.TraceContext{},
		),
	)),
}
