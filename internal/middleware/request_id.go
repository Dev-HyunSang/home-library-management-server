package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestIDMiddleware generates a UUID for each request and sets it
// in ctx.Locals and the X-Request-ID response header for log correlation.
func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Locals("requestID", requestID)
		c.Set("X-Request-ID", requestID)

		return c.Next()
	}
}
