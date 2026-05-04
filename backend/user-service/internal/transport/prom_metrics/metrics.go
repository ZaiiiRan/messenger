package prommetrics

import (
	"context"
	"net/http"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/config/settings"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	registry   *prometheus.Registry
	metricsSrv *http.Server
}

func New(cfg settings.MetricsServerSettings) *Server {
	registry := prometheus.NewRegistry()

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	metricsSrv := &http.Server{
		Addr:    cfg.Port,
		Handler: mux,
	}

	return &Server{
		registry:   registry,
		metricsSrv: metricsSrv,
	}
}

func (s *Server) Start() error {
	return s.metricsSrv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	if s.metricsSrv != nil {
		return s.metricsSrv.Shutdown(ctx)
	}
	return nil
}

func (s *Server) Registry() *prometheus.Registry {
	return s.registry
}
