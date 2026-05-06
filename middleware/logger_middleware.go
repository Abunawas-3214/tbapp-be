package middleware

import (
	"time"

	"tbapp-be/common/logger"

	"github.com/gofiber/fiber/v3"
)

// StructuredLogger returns a custom middleware that logs HTTP requests using zerolog for Fiber v3
func StructuredLogger() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		// Handle request
		err := c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get logger
		l := logger.GetLogger()

		// Determine log level and message
		status := c.Response().StatusCode()
		msg := "Request"

		event := l.Info()
		if err != nil || status >= 500 {
			event = l.Error()
			if err != nil {
				msg = err.Error()
			}
		} else if status >= 400 {
			event = l.Warn()
		}

		event.
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", status).
			Str("ip", c.IP()).
			Dur("latency", latency).
			Str("user-agent", string(c.Request().Header.UserAgent())).
			Msg(msg)

		return err
	}
}
