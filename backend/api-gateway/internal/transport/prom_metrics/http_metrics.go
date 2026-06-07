package prommetrics

import "github.com/prometheus/client_golang/prometheus"

// HTTPMetrics holds Prometheus metrics describing the gateway HTTP surface.
type HTTPMetrics struct {
	RequestsTotal    *prometheus.CounterVec
	RequestDuration  *prometheus.HistogramVec
	RequestsInFlight prometheus.Gauge
}

func NewHTTPMetrics(reg *prometheus.Registry) *HTTPMetrics {
	m := &HTTPMetrics{
		RequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gateway_http_requests_total",
				Help: "Total number of HTTP requests handled by the gateway.",
			},
			[]string{"method", "route", "status"},
		),
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "gateway_http_request_duration_seconds",
				Help:    "Duration of HTTP requests handled by the gateway in seconds.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "route", "status"},
		),
		RequestsInFlight: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "gateway_http_requests_in_flight",
				Help: "Current number of HTTP requests being served by the gateway.",
			},
		),
	}
	reg.MustRegister(m.RequestsTotal, m.RequestDuration, m.RequestsInFlight)
	return m
}
