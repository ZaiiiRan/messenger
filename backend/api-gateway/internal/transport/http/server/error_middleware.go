package httpserver

import "github.com/gofiber/fiber/v2"

// errorMappingMiddleware converts an error returned by an inner handler or
// middleware into a fully written HTTP response (correct status code + body)
// using the shared error handler, and then returns nil.
//
// It is registered as the innermost middleware (just outside recovery) so that
// by the time control bubbles back out to the logging and metrics middleware
// the response already carries the final, properly mapped status code instead
// of the default 200. The same handler is also installed as Fiber's top-level
// ErrorHandler as a safety net for the outer middleware themselves.
func errorMappingMiddleware(handler fiber.ErrorHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := c.Next(); err != nil {
			return handler(c, err)
		}
		return nil
	}
}
