package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/cmp-platform/backend/internal/metrics"
)

// PrometheusMetrics middleware records Prometheus metrics
func PrometheusMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		metrics.RecordHTTPRequest(
			c.Request.Method,
			c.FullPath(),
			c.Writer.Status(),
			duration,
		)
	}
}
