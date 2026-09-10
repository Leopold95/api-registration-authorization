package logging

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// HTTP invokes the configured error handler before recording the final status.
func HTTP() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		if err != nil {
			if handlerErr := c.App().ErrorHandler(c, err); handlerErr != nil {
				c.Status(fiber.StatusInternalServerError)
				err = errors.Join(err, handlerErr)
			}
		}
		level := zerolog.InfoLevel
		status := c.Response().StatusCode()
		if status >= 500 {
			level = zerolog.ErrorLevel
		} else if status >= 400 {
			level = zerolog.WarnLevel
		}
		log.WithLevel(level).Str("method", c.Method()).Str("path", c.Path()).
			Int("status", status).Float64("duration_ms", float64(time.Since(start))/float64(time.Millisecond)).
			Err(err).Msg("HTTP request completed")
		return nil
	}
}
