package metrics

import (
	"runtime"
	"strconv"
	"time"

	"api.us4ever/internal/logger"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"go.uber.org/zap"
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

	// 新增：自定义系统指标
	systemCPUUsage = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_cpu_usage",
			Help: "CPU usage percentage",
		},
		[]string{"cpu"},
	)

	systemMemoryUsage = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_memory_usage",
			Help: "Memory usage in bytes",
		},
		[]string{"type"},
	)

	systemDiskUsage = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_disk_usage",
			Help: "Disk usage in bytes",
		},
		[]string{"device", "mountpoint", "fstype", "mode"},
	)

	systemNetworkIO = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_network_io",
			Help: "Network IO in bytes",
		},
		[]string{"interface", "direction"},
	)

	// 新增：自定义运行时指标
	appGoroutines = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "app_goroutines",
			Help: "Number of goroutines",
		},
	)

	appCGOCalls = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "app_cgo_calls",
			Help: "Number of CGO calls",
		},
	)

	appMemStats = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "app_memstats",
			Help: "Memory statistics",
		},
		[]string{"stat"},
	)

	// 新增：业务指标
	indexSizeBytes = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "index_size_bytes",
			Help: "Size of the index in bytes",
		},
		[]string{"index"},
	)

	indexDocumentCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "index_document_count",
			Help: "Number of documents in the index",
		},
		[]string{"index"},
	)

	searchLatencyPercentile = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "search_latency_percentile",
			Help: "Search latency percentiles in seconds",
		},
		[]string{"type", "percentile"},
	)
)

// NewMiddleware creates a Fiber middleware for collecting HTTP metrics
func NewMiddleware() fiber.Handler {
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

// RecordIndexStats 新增：记录索引大小和文档数量
func RecordIndexStats(indexName string, sizeBytes int64, documentCount int) {
	indexSizeBytes.WithLabelValues(indexName).Set(float64(sizeBytes))
	indexDocumentCount.WithLabelValues(indexName).Set(float64(documentCount))
}

// RecordSearchLatencyPercentile 新增：记录搜索延迟百分位数
func RecordSearchLatencyPercentile(searchType string, percentile string, latency float64) {
	searchLatencyPercentile.WithLabelValues(searchType, percentile).Set(latency)
}

// Collector MetricsCollector provides methods to collect various application metrics
type Collector struct {
	logger       *logger.Logger
	registry     *prometheus.Registry
	collectEvery time.Duration
	stopCh       chan struct{}
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() (*Collector, error) {
	metricsLogger, err := logger.New("metrics")
	if err != nil {
		return nil, err
	}

	// 创建自定义注册表
	registry := prometheus.NewRegistry()

	// 注册标准 Go 运行时指标收集器
	registry.MustRegister(collectors.NewGoCollector())

	// 注册进程指标收集器
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	return &Collector{
		logger:       metricsLogger,
		registry:     registry,
		collectEvery: 30 * time.Second,
		stopCh:       make(chan struct{}),
	}, nil
}

// StartPeriodicCollection starts collecting metrics periodically
func (mc *Collector) StartPeriodicCollection() {
	ticker := time.NewTicker(mc.collectEvery)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				mc.collectAllMetrics()
			case <-mc.stopCh:
				mc.logger.Info("stopping metrics collection")
				return
			}
		}
	}()
}

// StopCollection stops the periodic metrics collection
func (mc *Collector) StopCollection() {
	close(mc.stopCh)
}

// collectAllMetrics collects all metrics
func (mc *Collector) collectAllMetrics() {
	start := time.Now()
	mc.logger.Debug("starting metrics collection")

	// 收集各类指标
	mc.collectRuntimeMetrics()
	mc.collectSystemMetrics()

	duration := time.Since(start)
	mc.logger.Debug("metrics collection completed",
		zap.Duration("duration", duration),
	)
}

// collectRuntimeMetrics collects Go runtime metrics
func (mc *Collector) collectRuntimeMetrics() {
	// 收集 goroutine 数量
	appGoroutines.Set(float64(runtime.NumGoroutine()))

	// 收集 CGO 调用次数
	appCGOCalls.Set(float64(runtime.NumCgoCall()))

	// 收集内存统计信息
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	appMemStats.WithLabelValues("alloc").Set(float64(memStats.Alloc))
	appMemStats.WithLabelValues("total_alloc").Set(float64(memStats.TotalAlloc))
	appMemStats.WithLabelValues("sys").Set(float64(memStats.Sys))
	appMemStats.WithLabelValues("lookups").Set(float64(memStats.Lookups))
	appMemStats.WithLabelValues("mallocs").Set(float64(memStats.Mallocs))
	appMemStats.WithLabelValues("frees").Set(float64(memStats.Frees))

	appMemStats.WithLabelValues("heap_alloc").Set(float64(memStats.HeapAlloc))
	appMemStats.WithLabelValues("heap_sys").Set(float64(memStats.HeapSys))
	appMemStats.WithLabelValues("heap_idle").Set(float64(memStats.HeapIdle))
	appMemStats.WithLabelValues("heap_inuse").Set(float64(memStats.HeapInuse))
	appMemStats.WithLabelValues("heap_released").Set(float64(memStats.HeapReleased))
	appMemStats.WithLabelValues("heap_objects").Set(float64(memStats.HeapObjects))

	appMemStats.WithLabelValues("stack_inuse").Set(float64(memStats.StackInuse))
	appMemStats.WithLabelValues("stack_sys").Set(float64(memStats.StackSys))
	appMemStats.WithLabelValues("mspan_inuse").Set(float64(memStats.MSpanInuse))
	appMemStats.WithLabelValues("mspan_sys").Set(float64(memStats.MSpanSys))
	appMemStats.WithLabelValues("mcache_inuse").Set(float64(memStats.MCacheInuse))
	appMemStats.WithLabelValues("mcache_sys").Set(float64(memStats.MCacheSys))

	appMemStats.WithLabelValues("gc_next").Set(float64(memStats.NextGC))
	appMemStats.WithLabelValues("gc_last").Set(float64(memStats.LastGC))
	appMemStats.WithLabelValues("gc_num").Set(float64(memStats.NumGC))
	appMemStats.WithLabelValues("gc_cpu_fraction").Set(memStats.GCCPUFraction)
}

// collectSystemMetrics collects system-level metrics
func (mc *Collector) collectSystemMetrics() {
	// 收集 CPU 使用率
	cpuPercent, err := cpu.Percent(0, true)
	if err != nil {
		mc.logger.Error("failed to collect CPU metrics", zap.Error(err))
	} else {
		for i, percent := range cpuPercent {
			systemCPUUsage.WithLabelValues(strconv.Itoa(i)).Set(percent)
		}
	}

	// 收集内存使用情况
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		mc.logger.Error("failed to collect memory metrics", zap.Error(err))
	} else {
		systemMemoryUsage.WithLabelValues("total").Set(float64(vmStat.Total))
		systemMemoryUsage.WithLabelValues("available").Set(float64(vmStat.Available))
		systemMemoryUsage.WithLabelValues("used").Set(float64(vmStat.Used))
		systemMemoryUsage.WithLabelValues("free").Set(float64(vmStat.Free))
		systemMemoryUsage.WithLabelValues("cached").Set(float64(vmStat.Cached))
		systemMemoryUsage.WithLabelValues("buffers").Set(float64(vmStat.Buffers))
	}

	// 收集磁盘使用情况
	partitions, err := disk.Partitions(false)
	if err != nil {
		mc.logger.Error("failed to collect disk partition metrics", zap.Error(err))
	} else {
		for _, partition := range partitions {
			usage, err := disk.Usage(partition.Mountpoint)
			if err != nil {
				mc.logger.Error("failed to collect disk usage metrics",
					zap.Error(err),
					zap.String("mountpoint", partition.Mountpoint),
				)
				continue
			}

			systemDiskUsage.WithLabelValues(
				partition.Device,
				partition.Mountpoint,
				partition.Fstype,
				"total",
			).Set(float64(usage.Total))

			systemDiskUsage.WithLabelValues(
				partition.Device,
				partition.Mountpoint,
				partition.Fstype,
				"free",
			).Set(float64(usage.Free))

			systemDiskUsage.WithLabelValues(
				partition.Device,
				partition.Mountpoint,
				partition.Fstype,
				"used",
			).Set(float64(usage.Used))
		}
	}

	// 收集网络 IO 统计
	netIOStats, err := net.IOCounters(true)
	if err != nil {
		mc.logger.Error("failed to collect network IO metrics", zap.Error(err))
	} else {
		for _, netStat := range netIOStats {
			systemNetworkIO.WithLabelValues(netStat.Name, "bytes_sent").Set(float64(netStat.BytesSent))
			systemNetworkIO.WithLabelValues(netStat.Name, "bytes_recv").Set(float64(netStat.BytesRecv))
			systemNetworkIO.WithLabelValues(netStat.Name, "packets_sent").Set(float64(netStat.PacketsSent))
			systemNetworkIO.WithLabelValues(netStat.Name, "packets_recv").Set(float64(netStat.PacketsRecv))
			systemNetworkIO.WithLabelValues(netStat.Name, "errin").Set(float64(netStat.Errin))
			systemNetworkIO.WithLabelValues(netStat.Name, "errout").Set(float64(netStat.Errout))
			systemNetworkIO.WithLabelValues(netStat.Name, "dropin").Set(float64(netStat.Dropin))
			systemNetworkIO.WithLabelValues(netStat.Name, "dropout").Set(float64(netStat.Dropout))
		}
	}
}

// GetMetricsHandler returns a Fiber handler for the /metrics endpoint
func GetMetricsHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		// 使用fiber的adaptor中间件将promhttp.Handler转换为fiber.Handler
		handler := adaptor.HTTPHandler(promhttp.Handler())
		return handler(c)
	}
}

// StartMetricsCollection 初始化并启动指标收集
func StartMetricsCollection() (*Collector, error) {
	collector, err := NewMetricsCollector()
	if err != nil {
		return nil, err
	}

	// 启动周期性收集
	collector.StartPeriodicCollection()

	return collector, nil
}
