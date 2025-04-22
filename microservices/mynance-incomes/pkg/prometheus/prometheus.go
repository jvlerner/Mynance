package prometheus

import (
	"database/sql"
	"log"
	"net/http"
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests, labeled by route, method, and status.",
		},
		[]string{"route", "method", "status"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds, labeled by route and method.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"route", "method"},
	)

	RequestsInProgress = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_in_progress",
			Help: "Number of HTTP requests currently in progress, labeled by route and method.",
		},
		[]string{"route", "method"},
	)

	ResponseSizeBytes = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Size of HTTP responses in bytes.",
			Buckets: prometheus.ExponentialBuckets(100, 2, 10),
		},
		[]string{"route", "method"},
	)

	ErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_errors_total",
			Help: "Total number of 5xx server errors.",
		},
		[]string{"route", "method"},
	)

	Goroutines = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "go_goroutines",
			Help: "Number of goroutines that currently exist.",
		},
		func() float64 {
			return float64(runtime.NumGoroutine())
		},
	)

	DBOpenConns       prometheus.GaugeFunc
	metricsServer     *http.Server
	dbGaugeRegistered bool
)

func SetDBForMonitoring(db *sql.DB) {
	DBOpenConns = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "db_open_connections",
			Help: "Number of open database connections.",
		},
		func() float64 {
			return float64(db.Stats().OpenConnections)
		},
	)
	prometheus.MustRegister(DBOpenConns)
	dbGaugeRegistered = true
}

func CloseDBForMonitoring() {
	if dbGaugeRegistered {
		prometheus.Unregister(DBOpenConns)
		dbGaugeRegistered = false
	}
}

func Init() {
	prometheus.MustRegister(
		RequestCounter,
		RequestDuration,
		RequestsInProgress,
		ResponseSizeBytes,
		ErrorCounter,
		Goroutines,
	)

	metricsServer = &http.Server{
		Addr:    ":2222",
		Handler: http.DefaultServeMux,
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("[INFO] Serving Prometheus metrics at :2222/metrics")
		if err := metricsServer.ListenAndServe(); err != nil {
			log.Printf("[ERROR] Failed to start /metrics endpoint: %v", err)
		}
	}()
}

// Close gracefully shuts down the /metrics server
func Close() {
	if metricsServer != nil {
		_ = metricsServer.Close() // graceful shutdown
	}
}
