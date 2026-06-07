package httpserver

import (
	"strconv"
	"time"

	prommetrics "github.com/ZaiiiRan/messenger/backend/api-gateway/internal/transport/prom_metrics"
	"github.com/gofiber/fiber/v2"
)

// metricsMiddleware records per-request Prometheus metrics for the gateway.
func metricsMiddleware(m *prommetrics.HTTPMetrics) fiber.Handler {
	return func(c *fiber.Ctx) error {
		m.RequestsInFlight.Inc()
		start := time.Now()

		err := c.Next()

		m.RequestsInFlight.Dec()

		route := c.Route().Path
		if route == "" {
			route = c.Path()
		}
		status := strconv.Itoa(c.Response().StatusCode())
		method := c.Method()

		m.RequestsTotal.WithLabelValues(method, route, status).Inc()
		m.RequestDuration.WithLabelValues(method, route, status).Observe(time.Since(start).Seconds())

		return err
	}
}
