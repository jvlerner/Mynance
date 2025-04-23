package prometheus

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	serviceName = getServiceName()

	RequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests, labeled by route, method, and status.",
		},
		[]string{"service", "route", "method", "status"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds, labeled by route and method.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "route", "method"},
	)

	RequestsInProgress = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_in_progress",
			Help: "Number of HTTP requests currently in progress, labeled by route and method.",
		},
		[]string{"service", "route", "method"},
	)

	ResponseSizeBytes = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Size of HTTP responses in bytes.",
			Buckets: prometheus.ExponentialBuckets(100, 2, 10),
		},
		[]string{"service", "route", "method"},
	)

	ErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_errors_total",
			Help: "Total number of 5xx server errors.",
		},
		[]string{"service", "route", "method"},
	)

	Goroutines = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "go_goroutines",
			Help: "Number of goroutines that currently exist.",
		},
		[]string{"service"},
	)

	dbOpenConns = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "db_open_connections",
			Help: "Number of open database connections.",
		},
		[]string{"db_name"},
	)
	dbGaugeRegistered = false

	metricsServer *http.Server
	initOnce      sync.Once
)

// getServiceName reads SERVICE_NAME or defaults to "unknown"
func getServiceName() string {
	if name := os.Getenv("SERVICE_NAME"); name != "" {
		return name
	}
	return "unknown"
}

// SetDBForMonitoring registers a labeled gauge for a given DB instance
func SetDBForMonitoring(db *sql.DB, name string) {
	if !dbGaugeRegistered {
		prometheus.MustRegister(dbOpenConns)
		dbGaugeRegistered = true
	}

	go func() {
		for {
			dbOpenConns.WithLabelValues(name).Set(float64(db.Stats().OpenConnections))
			time.Sleep(10 * time.Second)
		}
	}()
}

// RemoveDBFromMonitoring removes a DB gauge from Prometheus monitoring
func RemoveDBFromMonitoring(name string) {
	if dbGaugeRegistered {
		dbOpenConns.DeleteLabelValues(name)
	}
}

// Init starts Prometheus metrics collection and exposes /metrics
func Init() {
	initOnce.Do(func() {
		prometheus.MustRegister(
			RequestCounter,
			RequestDuration,
			RequestsInProgress,
			ResponseSizeBytes,
			ErrorCounter,
			Goroutines,
		)

		// Atualiza goroutines a cada 5 segundos
		go func() {
			ticker := time.NewTicker(5 * time.Second)
			for range ticker.C {
				Goroutines.WithLabelValues(serviceName).Set(float64(runtime.NumGoroutine()))
			}
		}()

		http.Handle("/metrics", promhttp.Handler())
		metricsServer = &http.Server{Addr: ":2222"}

		go func() {
			log.Println("[INFO] Serving Prometheus metrics at :2222/metrics")
			if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Printf("[ERROR] Failed to start /metrics endpoint: %v", err)
			}
		}()
	})
}

// Close gracefully shuts down the /metrics server
func Close() {
	if metricsServer != nil {
		_ = metricsServer.Close()
	}
}
