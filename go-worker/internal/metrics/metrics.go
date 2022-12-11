package metrics

import (
	"context"
	"fmt"
	"go-worker/internal/settings"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/metric"
)

var (
	Provider *metric.MeterProvider
	Meter    otelmetric.Meter
)

const (
	_responseTimeDuration = "response.time.duration"
	_successRequestCount  = "success.request"
	_errorRequestCount    = "error.request"
)

func StartupMetrics() {

	var err error
	Provider, err = NewMetrictProvider()

	if err != nil {
		log.Fatalf("Error when try create new metric provider: %s", err)
	}

	Meter = NewMeter(Provider)
	go StartMetricsEndpoint()
}

func NewMetrictProvider() (*metric.MeterProvider, error) {
	exporter, err := prometheus.New()
	if err != nil {
		return &metric.MeterProvider{}, err
	}
	return metric.NewMeterProvider(metric.WithReader(exporter)), nil
}

func NewMeter(provider *metric.MeterProvider) otelmetric.Meter {
	return provider.Meter("github.com/open-telemetry/opentelemetry-go/example/prometheus")
}

func StartMetricsEndpoint() {

	http.Handle("/", promhttp.Handler())

	timeOut := time.Duration(settings.Env.Metrics.Timeout) * time.Second
	port := fmt.Sprintf(":%d", settings.Env.Metrics.Port)

	server := &http.Server{
		Addr:              port,
		ReadHeaderTimeout: timeOut,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("error serving http: %v", err)
		return
	}

	log.Printf("serving metrics at %s/metrics", port)
}

func ResponseTimeMetric(duration time.Duration, attrs []attribute.KeyValue) {

	ctx := context.Background()
	timer, err := Meter.SyncInt64().Histogram(_responseTimeDuration, instrument.WithUnit(unit.Milliseconds))

	if err != nil {
		logGFLMetricError("ResponseTimeMetric", err)
		return
	}

	timer.Record(ctx, duration.Milliseconds(), attrs...)
}

func SuccessRequestCountMetric(attrs []attribute.KeyValue) {

	ctx := context.Background()
	success, err := Meter.SyncInt64().Counter(_successRequestCount)

	if err != nil {
		logGFLMetricError("SuccessRequestCountMetric", err)
		return
	}

	success.Add(ctx, 1, attrs...)
}

func ErrorRequestCountMetric(attrs []attribute.KeyValue) {

	ctx := context.Background()
	success, err := Meter.SyncInt64().Counter(_errorRequestCount)

	if err != nil {
		logGFLMetricError("ErrorRequestCountMetric", err)
		return
	}

	success.Add(ctx, 1, attrs...)
}

func BuildMetrictsAttrs(attrs ...attribute.KeyValue) []attribute.KeyValue {
	return append(attrs, attribute.Key("service_name").String(settings.Env.AppName))
}

func logGFLMetricError(method string, err error) {
	fields := log.Fields{"package": "metrics", "method": method}
	log.WithFields(fields).Error(err.Error())
}
