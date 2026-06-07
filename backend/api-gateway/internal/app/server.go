package app

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/api-gateway/internal/config"
	grpcclient "github.com/ZaiiiRan/messenger/backend/api-gateway/internal/transport/client/grpc"
	"github.com/ZaiiiRan/messenger/backend/api-gateway/internal/transport/http/handlers"
	httpserver "github.com/ZaiiiRan/messenger/backend/api-gateway/internal/transport/http/server"
	gatewayi18n "github.com/ZaiiiRan/messenger/backend/api-gateway/internal/transport/i18n"
	prommetrics "github.com/ZaiiiRan/messenger/backend/api-gateway/internal/transport/prom_metrics"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/logger"
	grpc_prom "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ServerApp struct {
	cfg *config.ServerConfig
	log *zap.SugaredLogger

	clients    *grpcclient.Clients
	httpServer *httpserver.Server

	metricsServer *prommetrics.Server

	isHTTPStarted    chan bool
	isMetricsStarted chan bool
}

func NewServerApp() (*ServerApp, error) {
	cfg, err := config.LoadServerConfig()
	if err != nil {
		return nil, err
	}

	log, err := logger.New()
	if err != nil {
		return nil, err
	}

	return &ServerApp{
		cfg:              cfg,
		log:              log,
		isHTTPStarted:    make(chan bool),
		isMetricsStarted: make(chan bool),
	}, nil
}

func (a *ServerApp) Run(ctx context.Context) error {
	a.initI18n()

	a.initMetricsServer()
	a.startMetricsServer()

	grpcClientMetrics := grpc_prom.NewClientMetrics()
	a.metricsServer.Registry().MustRegister(grpcClientMetrics)

	if err := a.initClients(ctx, grpcClientMetrics); err != nil {
		return err
	}

	a.initHTTPServer()
	a.startHTTPServer()

	<-a.isMetricsStarted
	<-a.isHTTPStarted
	a.log.Infow("app.started")
	return nil
}

func (a *ServerApp) Stop(ctx context.Context) {
	a.log.Infow("app.stopping")

	shCtx, cancel := context.WithTimeout(ctx, time.Duration(a.cfg.Shutdown.ShutdownTimeout)*time.Second)
	defer cancel()

	if a.httpServer != nil {
		if err := a.httpServer.Stop(shCtx); err != nil {
			a.log.Errorw("app.http_stop_failed", "err", err)
		}
	}
	if a.clients != nil {
		a.clients.Close()
	}
	if a.metricsServer != nil {
		_ = a.metricsServer.Stop(shCtx)
	}

	a.log.Infow("app.stopped")
}

func (a *ServerApp) initI18n() {
	gatewayi18n.Init()
}

func (a *ServerApp) initMetricsServer() {
	a.metricsServer = prommetrics.New(a.cfg.MetricsServer)
}

func (a *ServerApp) startMetricsServer() {
	go func() {
		a.log.Infow("app.metrics_serve_start", "port", a.cfg.MetricsServer.Port)
		a.isMetricsStarted <- true
		if err := a.metricsServer.Start(); err != nil {
			a.log.Fatalw("app.metrics_serve_error", "err", err)
		}
	}()
}

func (a *ServerApp) initClients(ctx context.Context, grpcClientMetrics *grpc_prom.ClientMetrics) error {
	unaryExtra := []grpc.UnaryClientInterceptor{
		grpcClientMetrics.UnaryClientInterceptor(),
	}

	clients, err := grpcclient.New(ctx, a.cfg, unaryExtra)
	if err != nil {
		a.log.Errorw("app.grpc_clients_init_failed", "err", err)
		return err
	}
	a.clients = clients

	a.log.Infow("app.grpc_clients_connected")
	return nil
}

func (a *ServerApp) initHTTPServer() {
	httpMetrics := prommetrics.NewHTTPMetrics(a.metricsServer.Registry())
	h := handlers.New(a.clients, a.cfg)
	a.httpServer = httpserver.New(a.cfg, h, httpMetrics, a.log)
}

func (a *ServerApp) startHTTPServer() {
	go func() {
		a.log.Infow("app.http_serve_start", "port", a.cfg.HTTPServer.Port)
		a.isHTTPStarted <- true
		if err := a.httpServer.Start(); err != nil {
			a.log.Fatalw("app.http_serve_error", "err", err)
		}
	}()
}
