package prommetrics

import "github.com/prometheus/client_golang/prometheus"

type WorkerMetrics struct {
	CyclesTotal         *prometheus.CounterVec
	ProcessedItemsTotal *prometheus.CounterVec
	CycleDuration       *prometheus.HistogramVec
}

func NewWorkerMetrics(reg *prometheus.Registry) *WorkerMetrics {
	m := &WorkerMetrics{
		CyclesTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "worker_cycles_total",
				Help: "Total number of worker cycles executed.",
			},
			[]string{"worker_type", "status"},
		),
		ProcessedItemsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "worker_processed_items_total",
				Help: "Total number of items processed by workers.",
			},
			[]string{"worker_type"},
		),
		CycleDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "worker_cycle_duration_seconds",
				Help:    "Duration of a single worker cycle in seconds.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"worker_type"},
		),
	}
	reg.MustRegister(m.CyclesTotal, m.ProcessedItemsTotal, m.CycleDuration)
	return m
}
