package metrics

import (
	"strconv"
	"time"

	"api.us4ever/internal/logger"
	"github.com/gofiber/fiber/v3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	httpRequestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP request size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "path"},
	)

	httpResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "path", "status"},
	)

	// Database metrics
	databaseConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_connections_active",
			Help: "Number of active database connections",
		},
	)

	databaseConnectionsIdle = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_connections_idle",
			Help: "Number of idle database connections",
		},
	)

	databaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	// Elasticsearch metrics
	elasticsearchRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "elasticsearch_requests_total",
			Help: "Total number of Elasticsearch requests",
		},
		[]string{"operation", "index", "status"},
	)

	elasticsearchRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "elasticsearch_request_duration_seconds",
			Help:    "Elasticsearch request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "index"},
	)

	// Application metrics
	taskExecutionsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "task_executions_total",
			Help: "Total number of task executions",
		},
		[]string{"task", "status"},
	)

	taskExecutionDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "task_execution_duration_seconds",
			Help:    "Task execution duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"task"},
	)

	// Search metrics
	searchRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "search_requests_total",
			Help: "Total number of search requests",
		},
		[]string{"type", "status"},
	)

	searchResultsCount = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "search_results_count",
			Help:    "Number of search results returned",
			Buckets: []float64{0, 1, 5, 10, 25, 50, 100, 250, 500, 1000},
		},
		[]string{"type"},
	)
)

// MetricsMiddleware creates a Fiber middleware for collecting HTTP metrics
func MetricsMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		// Record request size
		if c.Request().Header.ContentLength() > 0 {
			httpRequestSize.WithLabelValues(
				c.Method(),
				c.Path(),
			).Observe(float64(c.Request().Header.ContentLength()))
		}

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Get status code
		status := strconv.Itoa(c.Response().StatusCode())

		// Record metrics
		httpRequestsTotal.WithLabelValues(
			c.Method(),
			c.Path(),
			status,
		).Inc()

		httpRequestDuration.WithLabelValues(
			c.Method(),
			c.Path(),
			status,
		).Observe(duration)

		// Record response size
		responseSize := len(c.Response().Body())
		if responseSize > 0 {
			httpResponseSize.WithLabelValues(
				c.Method(),
				c.Path(),
				status,
			).Observe(float64(responseSize))
		}

		return err
	}
}

// RecordDatabaseConnection records database connection metrics
func RecordDatabaseConnection(active, idle int) {
	databaseConnectionsActive.Set(float64(active))
	databaseConnectionsIdle.Set(float64(idle))
}

// RecordDatabaseQuery records database query metrics
func RecordDatabaseQuery(operation string, duration time.Duration) {
	databaseQueryDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

// RecordElasticsearchRequest records Elasticsearch request metrics
func RecordElasticsearchRequest(operation, index, status string, duration time.Duration) {
	elasticsearchRequestsTotal.WithLabelValues(operation, index, status).Inc()
	elasticsearchRequestDuration.WithLabelValues(operation, index).Observe(duration.Seconds())
}

// RecordTaskExecution records task execution metrics
func RecordTaskExecution(taskName, status string, duration time.Duration) {
	taskExecutionsTotal.WithLabelValues(taskName, status).Inc()
	taskExecutionDuration.WithLabelValues(taskName).Observe(duration.Seconds())
}

// RecordSearchRequest records search request metrics
func RecordSearchRequest(searchType, status string, resultCount int) {
	searchRequestsTotal.WithLabelValues(searchType, status).Inc()
	searchResultsCount.WithLabelValues(searchType).Observe(float64(resultCount))
}

// MetricsCollector provides methods to collect various application metrics
type MetricsCollector struct {
	logger *logger.Logger
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() (*MetricsCollector, error) {
	metricsLogger, err := logger.New("metrics")
	if err != nil {
		return nil, err
	}

	return &MetricsCollector{
		logger: metricsLogger,
	}, nil
}

// StartPeriodicCollection starts collecting metrics periodically
func (mc *MetricsCollector) StartPeriodicCollection() {
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			mc.collectSystemMetrics()
		}
	}()
}

// collectSystemMetrics collects system-level metrics
func (mc *MetricsCollector) collectSystemMetrics() {
	// This is a placeholder - in a real implementation, you would collect
	// actual system metrics like memory usage, CPU usage, etc.
	mc.logger.Debug("collecting system metrics")
}

// GetMetricsHandler returns a Fiber handler for the /metrics endpoint
func GetMetricsHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		// In a real implementation, you would use promhttp.Handler()
		// For now, return a simple response
		return c.SendString("# Metrics endpoint - integrate with Prometheus\n")
	}
}
