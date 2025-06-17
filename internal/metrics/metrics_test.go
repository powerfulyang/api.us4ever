package metrics

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestRuntimeMetricsCollection(t *testing.T) {
	// 创建一个新的收集器
	collector, err := NewMetricsCollector()
	if err != nil {
		t.Fatalf("Failed to create metrics collector: %v", err)
	}

	// 手动触发收集
	collector.collectRuntimeMetrics()

	// 验证 goroutine 指标是否被收集
	if testutil.ToFloat64(appGoroutines) <= 0 {
		t.Errorf("Expected appGoroutines to be greater than 0, got %v", testutil.ToFloat64(appGoroutines))
	}

	// 验证 CGO 调用指标是否被收集
	if testutil.ToFloat64(appCGOCalls) < 0 {
		t.Errorf("Expected appCGOCalls to be non-negative, got %v", testutil.ToFloat64(appCGOCalls))
	}

	// 验证内存统计信息是否被收集
	alloc := testutil.ToFloat64(appMemStats.WithLabelValues("alloc"))
	if alloc <= 0 {
		t.Errorf("Expected appMemStats[alloc] to be greater than 0, got %v", alloc)
	}
}

func TestSystemMetricsCollection(t *testing.T) {
	// 创建一个新的收集器
	collector, err := NewMetricsCollector()
	if err != nil {
		t.Fatalf("Failed to create metrics collector: %v", err)
	}

	// 手动触发收集
	collector.collectSystemMetrics()

	// 注意：由于系统指标收集可能会因权限或平台限制而失败，
	// 我们不能严格断言指标值，只能验证收集过程不会崩溃
}

func TestPeriodicCollection(t *testing.T) {
	// 创建一个新的收集器，使用较短的收集间隔
	collector, err := NewMetricsCollector()
	if err != nil {
		t.Fatalf("Failed to create metrics collector: %v", err)
	}
	collector.collectEvery = 100 * time.Millisecond

	// 启动周期性收集
	collector.StartPeriodicCollection()

	// 等待足够的时间，确保至少收集了一次
	time.Sleep(200 * time.Millisecond)

	// 停止收集
	collector.StopCollection()

	// 由于收集是异步的，我们无法直接验证结果
	// 但至少可以确保启动和停止过程不会崩溃
}

func TestHTTPMetricsRecording(t *testing.T) {
	// 重置指标
	httpRequestsTotal.Reset()
	httpRequestDuration.Reset()

	// 记录一个请求
	httpRequestsTotal.WithLabelValues("GET", "/test", "200").Inc()
	httpRequestDuration.WithLabelValues("GET", "/test", "200").Observe(0.1)

	// 验证指标是否被正确记录
	count := testutil.ToFloat64(httpRequestsTotal.WithLabelValues("GET", "/test", "200"))
	if count != 1 {
		t.Errorf("Expected httpRequestsTotal to be 1, got %v", count)
	}
}

func TestDatabaseMetricsRecording(t *testing.T) {
	// 重置指标
	databaseConnectionsActive.Set(0)
	databaseConnectionsIdle.Set(0)
	databaseQueryDuration.Reset()

	// 记录连接数
	RecordDatabaseConnection(5, 10)

	// 记录查询耗时
	RecordDatabaseQuery("select", 100*time.Millisecond)

	// 验证指标是否被正确记录
	active := testutil.ToFloat64(databaseConnectionsActive)
	if active != 5 {
		t.Errorf("Expected databaseConnectionsActive to be 5, got %v", active)
	}

	idle := testutil.ToFloat64(databaseConnectionsIdle)
	if idle != 10 {
		t.Errorf("Expected databaseConnectionsIdle to be 10, got %v", idle)
	}
}

func TestElasticsearchMetricsRecording(t *testing.T) {
	// 重置指标
	elasticsearchRequestsTotal.Reset()
	elasticsearchRequestDuration.Reset()

	// 记录请求
	RecordElasticsearchRequest("search", "test-index", "200", 50*time.Millisecond)

	// 验证指标是否被正确记录
	count := testutil.ToFloat64(elasticsearchRequestsTotal.WithLabelValues("search", "test-index", "200"))
	if count != 1 {
		t.Errorf("Expected elasticsearchRequestsTotal to be 1, got %v", count)
	}
}

func TestSearchMetricsRecording(t *testing.T) {
	// 重置指标
	searchRequestsTotal.Reset()
	searchResultsCount.Reset()

	// 记录搜索请求
	RecordSearchRequest("keeps", "200", 15)

	// 验证指标是否被正确记录
	count := testutil.ToFloat64(searchRequestsTotal.WithLabelValues("keeps", "200"))
	if count != 1 {
		t.Errorf("Expected searchRequestsTotal to be 1, got %v", count)
	}
}

func TestIndexStatsRecording(t *testing.T) {
	// 重置指标
	indexSizeBytes.Reset()
	indexDocumentCount.Reset()

	// 记录索引统计信息
	RecordIndexStats("test-index", 1024*1024, 100)

	// 验证指标是否被正确记录
	size := testutil.ToFloat64(indexSizeBytes.WithLabelValues("test-index"))
	if size != 1024*1024 {
		t.Errorf("Expected indexSizeBytes to be %v, got %v", 1024*1024, size)
	}

	count := testutil.ToFloat64(indexDocumentCount.WithLabelValues("test-index"))
	if count != 100 {
		t.Errorf("Expected indexDocumentCount to be 100, got %v", count)
	}
}

func TestSearchLatencyPercentileRecording(t *testing.T) {
	// 重置指标
	searchLatencyPercentile.Reset()

	// 记录搜索延迟百分位数
	RecordSearchLatencyPercentile("keeps", "p99", 0.5)

	// 验证指标是否被正确记录
	latency := testutil.ToFloat64(searchLatencyPercentile.WithLabelValues("keeps", "p99"))
	if latency != 0.5 {
		t.Errorf("Expected searchLatencyPercentile to be 0.5, got %v", latency)
	}
}
