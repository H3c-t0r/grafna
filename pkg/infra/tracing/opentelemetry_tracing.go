package tracing

import (
	"context"
	"net/http"
	"time"

	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/infra/log/level"
	"github.com/grafana/grafana/pkg/setting"
	"go.etcd.io/etcd/version"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	trace "go.opentelemetry.io/otel/trace"
)

type Tracer interface {
	Run(context.Context) error
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, Span)
	Inject(context.Context, http.Header, Span)
}

type Span interface {
	End()
	SetAttributes(key string, value interface{}, kv attribute.KeyValue)
	SetName(name string)
	SetStatus(code codes.Code, description string)
	RecordError(err error, options ...trace.EventOption)
	AddEvents(keys []string, values []EventValue)
}

type Opentelemetry struct {
	enabledJaeger bool
	enabledOTLP   bool
	address       string
	log           log.Logger

	tracerProvider *tracesdk.TracerProvider
	tracer         trace.Tracer

	Cfg *setting.Cfg
}

type OpentelemetrySpan struct {
	span trace.Span
}

type EventValue struct {
	Str string
	Num int64
}

type otelErrHandler func(err error)

func (o otelErrHandler) Handle(err error) {
	o(err)
}

func (ots *Opentelemetry) parseSettingsOpentelemetry() error {
	section, err := ots.Cfg.Raw.GetSection("tracing.opentelemetry.jaeger")
	if err != nil {
		return err
	}

	ots.address = section.Key("address").MustString("")
	if ots.address != "" {
		ots.enabledJaeger = true
		return nil
	}

	section, err = ots.Cfg.Raw.GetSection("tracing.opentelemetry.otlp")
	if err != nil {
		return err
	}

	ots.address = section.Key("address").MustString("")
	if ots.address != "" {
		ots.enabledOTLP = true
	}

	return nil
}

func (ots *Opentelemetry) initJaegerTracerProvider() (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(ots.address)))
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("grafana"),
			attribute.String("environment", "production"),
		)),
	)

	return tp, nil
}

func (ots *Opentelemetry) initOTLPTracerProvider() (*tracesdk.TracerProvider, error) {
	client := otlptracegrpc.NewClient(otlptracegrpc.WithEndpoint(ots.address), otlptracegrpc.WithInsecure())
	exp, err := otlptrace.New(context.Background(), client)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("grafana"),
			semconv.ServiceVersionKey.String(version.Version),
		),
		resource.WithProcessRuntimeDescription(),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithSampler(tracesdk.ParentBased(
			tracesdk.AlwaysSample(),
		)),
		tracesdk.WithResource(res),
	)
	return tp, nil
}

func (ots *Opentelemetry) initOpentelemetryTracer() error {
	var tp *tracesdk.TracerProvider
	var err error
	if ots.enabledJaeger {
		tp, err = ots.initJaegerTracerProvider()
		if err != nil {
			return err
		}
	}

	if ots.enabledOTLP {
		tp, err = ots.initOTLPTracerProvider()
		if err != nil {
			return err
		}
	}
	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it
	// only if tracing is enabled
	if ots.enabledJaeger || ots.enabledOTLP {
		otel.SetTracerProvider(tp)
	}

	ots.tracerProvider = tp
	ots.tracer = otel.GetTracerProvider().Tracer("component-main")

	return nil
}

func (ots *Opentelemetry) Run(ctx context.Context) error {
	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetErrorHandler(otelErrHandler(func(err error) {
		// nolint:gosec
		level.Error(ots.log).Log("msg", "OpenTelemetry handler returned an error", "err", err)
	}))
	<-ctx.Done()

	ots.log.Info("Closing tracing")
	if ots.tracerProvider == nil {
		return nil
	}
	ctxShutdown, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	if err := ots.tracerProvider.Shutdown(ctxShutdown); err != nil {
		return err
	}

	return nil
}

func (ots *Opentelemetry) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, Span) {
	ctx, span := ots.tracer.Start(ctx, spanName)
	opentelemetrySpan := OpentelemetrySpan{
		span: span,
	}
	ctx = context.WithValue(ctx, traceKey{}, traceValue{span.SpanContext().TraceID().String(), span.SpanContext().IsSampled()})
	return ctx, opentelemetrySpan
}

func (ots *Opentelemetry) Inject(ctx context.Context, header http.Header, _ Span) {
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(header))
}

func (s OpentelemetrySpan) End() {
	s.span.End()
}

func (s OpentelemetrySpan) SetAttributes(key string, value interface{}, kv attribute.KeyValue) {
	s.span.SetAttributes(kv)
}

func (s OpentelemetrySpan) SetName(name string) {
	s.span.SetName(name)
}

func (s OpentelemetrySpan) SetStatus(code codes.Code, description string) {
	s.span.SetStatus(code, description)
}

func (s OpentelemetrySpan) RecordError(err error, options ...trace.EventOption) {
	for _, o := range options {
		s.span.RecordError(err, o)
	}
}

func (s OpentelemetrySpan) AddEvents(keys []string, values []EventValue) {
	for i, v := range values {
		if v.Num != 0 {
			s.span.AddEvent(keys[i], trace.WithAttributes(attribute.Key(keys[i]).String(v.Str)))
		}
		if v.Str != "" {
			s.span.AddEvent(keys[i], trace.WithAttributes(attribute.Key(keys[i]).Int64(v.Num)))
		}
	}
}
