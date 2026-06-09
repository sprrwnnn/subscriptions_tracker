package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		start := time.Now()

		log.WithFields(logrus.Fields{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"client_ip":  c.ClientIP(),
		}).Info("Request started")

		c.Next()

		latency := time.Since(start)

		log.WithFields(logrus.Fields{
			"request_id": requestID,
			"status":     c.Writer.Status(),
			"latency_ms": latency.Milliseconds(),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
		}).Info("Request completed")
	}
}
