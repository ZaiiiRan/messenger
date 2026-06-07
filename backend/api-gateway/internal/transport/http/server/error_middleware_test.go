package httpserver

import (
	"net/http/httptest"
	"testing"

	"github.com/ZaiiiRan/messenger/backend/api-gateway/internal/transport/i18n"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/errors/commonerror"
	commonmiddleware "github.com/ZaiiiRan/messenger/backend/go-common/pkg/middleware/http/server"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func testLogger() *zap.SugaredLogger {
	return zap.NewNop().Sugar()
}

// buildTestApp wires the same middleware ordering used by the real gateway and
// captures the status code observed by an outer (log/metrics-like) middleware.
func buildTestApp(observed *int) *fiber.App {
	i18n.Init()

	app := fiber.New(fiber.Config{
		ErrorHandler: commonmiddleware.NewErrorHandler(),
	})

	// Outer middleware standing in for log / metrics: it reads the status only
	// after the inner chain (including error mapping) has returned.
	app.Use(func(c *fiber.Ctx) error {
		err := c.Next()
		*observed = c.Response().StatusCode()
		return err
	})
	app.Use(commonmiddleware.I18nMiddleware(i18n.NewLocalizer))
	app.Use(errorMappingMiddleware(commonmiddleware.NewErrorHandler()))
	app.Use(commonmiddleware.RecoveryMiddleware(testLogger()))

	return app
}

func TestErrorMapping(t *testing.T) {
	tests := []struct {
		name    string
		handler fiber.Handler
		want    int
	}{
		{
			name:    "grpc not found -> 404",
			handler: func(c *fiber.Ctx) error { return status.Error(codes.NotFound, "missing") },
			want:    fiber.StatusNotFound,
		},
		{
			name:    "grpc permission denied -> 403",
			handler: func(c *fiber.Ctx) error { return status.Error(codes.PermissionDenied, "denied") },
			want:    fiber.StatusForbidden,
		},
		{
			name:    "grpc invalid argument -> 400",
			handler: func(c *fiber.Ctx) error { return status.Error(codes.InvalidArgument, "bad") },
			want:    fiber.StatusBadRequest,
		},
		{
			name:    "grpc unauthenticated -> 401",
			handler: func(c *fiber.Ctx) error { return status.Error(codes.Unauthenticated, "no") },
			want:    fiber.StatusUnauthorized,
		},
		{
			name:    "common unauthorized -> 401",
			handler: func(c *fiber.Ctx) error { return commonerror.ErrUnauthorized },
			want:    fiber.StatusUnauthorized,
		},
		{
			name:    "fiber bad request -> 400",
			handler: func(c *fiber.Ctx) error { return fiber.NewError(fiber.StatusBadRequest, "bad body") },
			want:    fiber.StatusBadRequest,
		},
		{
			name:    "panic -> 500",
			handler: func(c *fiber.Ctx) error { panic("boom") },
			want:    fiber.StatusInternalServerError,
		},
		{
			name:    "success -> 200",
			handler: func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"ok": true}) },
			want:    fiber.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var observed int
			app := buildTestApp(&observed)
			app.Get("/", tc.handler)

			resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}

			if resp.StatusCode != tc.want {
				t.Errorf("response status = %d, want %d", resp.StatusCode, tc.want)
			}
			if observed != tc.want {
				t.Errorf("status observed by outer middleware = %d, want %d (log/metrics would be wrong)", observed, tc.want)
			}
		})
	}
}
