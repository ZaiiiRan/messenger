package httpserver

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/api-gateway/docs"
	"github.com/ZaiiiRan/messenger/backend/api-gateway/internal/config"
	"github.com/ZaiiiRan/messenger/backend/api-gateway/internal/transport/http/handlers"
	gatewayi18n "github.com/ZaiiiRan/messenger/backend/api-gateway/internal/transport/i18n"
	prommetrics "github.com/ZaiiiRan/messenger/backend/api-gateway/internal/transport/prom_metrics"
	commonmiddleware "github.com/ZaiiiRan/messenger/backend/go-common/pkg/middleware/http/server"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"go.uber.org/zap"
)

type Server struct {
	cfg  *config.ServerConfig
	app  *fiber.App
	log  *zap.SugaredLogger
	addr string
}

func New(
	cfg *config.ServerConfig,
	h *handlers.Handlers,
	httpMetrics *prommetrics.HTTPMetrics,
	log *zap.SugaredLogger,
) *Server {
	app := fiber.New(fiber.Config{
		AppName:               "api-gateway",
		ErrorHandler:          commonmiddleware.NewErrorHandler(),
		ReadTimeout:           time.Duration(cfg.HTTPServer.ReadTimeout) * time.Second,
		WriteTimeout:          time.Duration(cfg.HTTPServer.WriteTimeout) * time.Second,
		IdleTimeout:           time.Duration(cfg.HTTPServer.IdleTimeout) * time.Second,
		BodyLimit:             cfg.HTTPServer.BodyLimit,
		DisableStartupMessage: true,
		ProxyHeader:           cfg.HTTPServer.TrustedProxyHeader,
	})

	s := &Server{
		cfg:  cfg,
		app:  app,
		log:  log,
		addr: cfg.HTTPServer.Port,
	}

	s.registerMiddleware(httpMetrics)
	s.registerSystemRoutes()
	s.registerRoutes(h)

	return s
}

func (s *Server) registerMiddleware(httpMetrics *prommetrics.HTTPMetrics) {
	if s.cfg.CORS.Enabled {
		s.app.Use(cors.New(cors.Config{
			AllowOrigins:     s.cfg.CORS.AllowOrigins,
			AllowMethods:     s.cfg.CORS.AllowMethods,
			AllowHeaders:     s.cfg.CORS.AllowHeaders,
			AllowCredentials: s.cfg.CORS.AllowCredentials,
		}))
	}

	// Middleware order is outer -> inner. Logging and metrics are kept outer so
	// they observe the final response status, which the error-mapping
	// middleware sets before control returns to them. Recovery is innermost so
	// a panic in a handler is turned into an error that the error-mapping
	// middleware can map to a proper status code.
	s.app.Use(commonmiddleware.RequestIdMiddleware())
	s.app.Use(commonmiddleware.LogMiddleware(s.log))
	s.app.Use(metricsMiddleware(httpMetrics))
	s.app.Use(commonmiddleware.I18nMiddleware(gatewayi18n.NewLocalizer))
	s.app.Use(errorMappingMiddleware(commonmiddleware.NewErrorHandler()))
	s.app.Use(commonmiddleware.RecoveryMiddleware(s.log))
}

func (s *Server) registerSystemRoutes() {
	s.app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	if s.cfg.Swagger.Enabled {
		s.app.Get("/openapi.json", func(c *fiber.Ctx) error {
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			return c.Send(docs.OpenAPISpec)
		})
		s.app.Get(s.cfg.Swagger.BasePath+"/*", swagger.New(swagger.Config{
			URL: "/openapi.json",
		}))
	}
}

func (s *Server) Start() error {
	return s.app.Listen(s.addr)
}

func (s *Server) Stop(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}
