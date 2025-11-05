package middleware

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func HttpLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()
		err := c.Next()
		duration := time.Since(startTime)
		statusCode := c.Response().StatusCode()

		logger := log.Info()
		if statusCode >= fiber.StatusBadRequest {
			logger = log.Error().Bytes("body", c.Response().Body())
		}

		logger.
			Str("method", c.Method()).
			Str("path", c.OriginalURL()).
			Int("status_code", statusCode).
			Str("status_text", http.StatusText(statusCode)).
			Dur("duration", duration).
			Msg("HTTP")

		return err
	}
}
