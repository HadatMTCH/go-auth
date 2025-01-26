package middleware

import (
	"base-api/infra/log"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/labstack/echo/v4"
)

func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Generate trace ID
			traceID := generateTraceID()
			c.Set("trace_id", traceID)

			// Process request
			err := next(c)

			// Collect log data
			latency := time.Since(start)
			status := c.Response().Status

			entry := log.FromContext(c).WithFields(logrus.Fields{
				"duration": latency.String(),
			})

			switch {
			case status >= 500:
				entry.WithField("error", err).Error("Server error")
			case status >= 400:
				entry.Warn("Client error")
			default:
				entry.Info("Request processed")
			}

			return err
		}
	}
}

func generateTraceID() string {
	// Implement your preferred trace ID generation
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
