package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jvlerner/my-finance-api/pkg/prometheus"
)

func Prometheus() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		route := c.FullPath()
		if route == "" {
			route = "undefined"
		}
		method := c.Request.Method

		prometheus.RequestsInProgress.WithLabelValues(route, method).Inc()

		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		duration := time.Since(start).Seconds()
		responseSize := float64(c.Writer.Size())

		prometheus.RequestCounter.WithLabelValues(route, method, status).Inc()
		prometheus.RequestDuration.WithLabelValues(route, method).Observe(duration)
		prometheus.RequestsInProgress.WithLabelValues(route, method).Dec()
		prometheus.ResponseSizeBytes.WithLabelValues(route, method).Observe(responseSize)

		if c.Writer.Status() >= 500 {
			prometheus.ErrorCounter.WithLabelValues(route, method).Inc()
		}
	}
}
