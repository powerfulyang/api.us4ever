# 代码改进集成示例

## 如何集成这些改进到现有代码中

### 1. 在路由中添加验证和速率限制

```go
// internal/server/routes.go - 更新搜索路由
func (s *FiberServer) setupRoutes() {
    // 添加全局中间件
    s.App.Use(metrics.MetricsMiddleware())
    
    // API路由组
    api := s.App.Group("/internal")
    
    // 为搜索端点添加专门的速率限制
    searchGroup := api.Group("/search")
    searchGroup.Use(middleware.NewSearchRateLimiter())
    
    // 搜索路由
    searchGroup.Get("/keeps", s.searchKeepsHandler)
    searchGroup.Get("/moments", s.searchMomentsHandler)
    
    // 其他路由使用通用速率限制
    api.Use(middleware.NewIPRateLimiter(30)) // 30 requests per second
    api.Get("/health", s.healthHandler)
}

// 更新搜索处理器以使用验证
func (s *FiberServer) searchKeepsHandler(c *fiber.Ctx) error {
    start := time.Now()
    
    // 解析查询参数
    query := c.Query("q", "")
    limit := c.QueryInt("limit", 10)
    offset := c.QueryInt("offset", 0)
    
    // 创建搜索请求
    searchReq := &validator.SearchRequest{
        Query:  query,
        Limit:  limit,
        Offset: offset,
    }
    
    // 验证请求
    if err := validator.ValidateSearchRequest(searchReq); err != nil {
        metrics.RecordSearchRequest("keeps", "validation_error", 0)
        return c.Status(400).JSON(fiber.Map{
            "error": fiber.Map{
                "type":    "ValidationError",
                "message": err.Error(),
                "code":    400,
            },
        })
    }
    
    // 清理查询
    sanitizedQuery, err := validator.ValidateAndSanitizeQuery(searchReq.Query)
    if err != nil {
        metrics.RecordSearchRequest("keeps", "sanitization_error", 0)
        return c.Status(400).JSON(fiber.Map{
            "error": fiber.Map{
                "type":    "SanitizationError",
                "message": err.Error(),
                "code":    400,
            },
        })
    }
    
    // 执行搜索
    ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
    defer cancel()
    
    results, err := es.SearchKeeps(ctx, s.EsClient, "api-keeps", sanitizedQuery)
    if err != nil {
        duration := time.Since(start)
        metrics.RecordSearchRequest("keeps", "search_error", 0)
        metrics.RecordElasticsearchRequest("search", "api-keeps", "error", duration)
        
        return c.Status(500).JSON(fiber.Map{
            "error": fiber.Map{
                "type":    "SearchError",
                "message": "Search operation failed",
                "code":    500,
            },
        })
    }
    
    // 记录成功的搜索指标
    duration := time.Since(start)
    resultCount := len(results.Hits.Hits)
    metrics.RecordSearchRequest("keeps", "success", resultCount)
    metrics.RecordElasticsearchRequest("search", "api-keeps", "success", duration)
    
    return c.JSON(fiber.Map{
        "query":    sanitizedQuery,
        "total":    results.Hits.Total.Value,
        "results":  results.Hits.Hits,
        "duration": duration.Milliseconds(),
    })
}
```

### 2. 在主函数中初始化改进

```go
// cmd/api/main.go - 更新主函数
func main() {
    // 初始化配置
    appConfig := config.GetAppConfig()
    if appConfig == nil {
        logger.Fatal("failed to load application config")
    }
    
    // 验证配置
    if err := appConfig.Validate(); err != nil {
        logger.Fatal("config validation failed", logger.Fields{
            "error": err.Error(),
        })
    }
    
    // 初始化指标收集器
    metricsCollector, err := metrics.NewMetricsCollector()
    if err != nil {
        logger.Fatal("failed to initialize metrics collector", logger.Fields{
            "error": err.Error(),
        })
    }
    metricsCollector.StartPeriodicCollection()
    
    // 创建服务器
    server, err := server.NewFiberServer(appConfig)
    if err != nil {
        logger.Fatal("failed to create server", logger.Fields{
            "error": err.Error(),
        })
    }
    
    // 添加指标端点
    server.App.Get("/metrics", metrics.GetMetricsHandler())
    
    // 启动服务器
    logger.Info("starting server", logger.Fields{
        "port": appConfig.Server.Port,
        "env":  appConfig.Environment,
    })
    
    if err := server.Listen(fmt.Sprintf(":%d", appConfig.Server.Port)); err != nil {
        logger.Fatal("server failed to start", logger.Fields{
            "error": err.Error(),
        })
    }
}
```

### 3. 在ES操作中添加指标记录

```go
// internal/es/search.go - 更新搜索函数
func SearchKeeps(ctx context.Context, client *elasticsearch.Client, indexName, query string) (*SearchResponse, error) {
    start := time.Now()
    
    // 构建搜索请求
    searchBody := map[string]interface{}{
        "query": map[string]interface{}{
            "multi_match": map[string]interface{}{
                "query":  query,
                "fields": []string{"title^2", "content", "summary"},
            },
        },
        "highlight": map[string]interface{}{
            "fields": map[string]interface{}{
                "title":   map[string]interface{}{},
                "content": map[string]interface{}{},
                "summary": map[string]interface{}{},
            },
        },
    }
    
    // 执行搜索
    res, err := client.Search(
        client.Search.WithContext(ctx),
        client.Search.WithIndex(indexName),
        client.Search.WithBody(esutil.NewJSONReader(searchBody)),
        client.Search.WithTrackTotalHits(true),
    )
    
    duration := time.Since(start)
    
    if err != nil {
        metrics.RecordElasticsearchRequest("search", indexName, "error", duration)
        return nil, fmt.Errorf("search request failed: %w", err)
    }
    defer res.Body.Close()
    
    if res.IsError() {
        metrics.RecordElasticsearchRequest("search", indexName, "error", duration)
        return nil, fmt.Errorf("search returned error: %s", res.Status())
    }
    
    // 解析响应
    var searchResponse SearchResponse
    if err := json.NewDecoder(res.Body).Decode(&searchResponse); err != nil {
        metrics.RecordElasticsearchRequest("search", indexName, "decode_error", duration)
        return nil, fmt.Errorf("failed to decode search response: %w", err)
    }
    
    metrics.RecordElasticsearchRequest("search", indexName, "success", duration)
    return &searchResponse, nil
}
```

### 4. 在任务执行中添加指标

```go
// internal/task/scheduler.go - 更新任务执行
func (s *Scheduler) executeTask(taskName string, taskFunc func() error) {
    start := time.Now()
    
    logger.Info("starting task execution", logger.Fields{
        "task": taskName,
    })
    
    err := taskFunc()
    duration := time.Since(start)
    
    if err != nil {
        metrics.RecordTaskExecution(taskName, "error", duration)
        logger.Error("task execution failed", logger.Fields{
            "task":     taskName,
            "duration": duration.String(),
            "error":    err.Error(),
        })
    } else {
        metrics.RecordTaskExecution(taskName, "success", duration)
        logger.Info("task execution completed", logger.Fields{
            "task":     taskName,
            "duration": duration.String(),
        })
    }
}
```

### 5. 配置验证示例

```go
// internal/config/validation.go
func (c *AppConfig) Validate() error {
    var errs []error
    
    // 验证应用名称
    if c.AppName == "" {
        errs = append(errs, errors.New("app_name is required"))
    }
    
    // 验证服务器配置
    if c.Server.Port <= 0 || c.Server.Port > 65535 {
        errs = append(errs, errors.New("server port must be between 1 and 65535"))
    }
    
    // 验证数据库配置
    if c.Database.Host == "" {
        errs = append(errs, errors.New("database host is required"))
    }
    if c.Database.Port <= 0 {
        errs = append(errs, errors.New("database port must be positive"))
    }
    if c.Database.Database == "" {
        errs = append(errs, errors.New("database name is required"))
    }
    
    // 验证ES配置
    if len(c.ES.Addresses) == 0 {
        errs = append(errs, errors.New("elasticsearch addresses are required"))
    }
    
    if len(errs) > 0 {
        return fmt.Errorf("config validation failed: %v", errs)
    }
    
    return nil
}
```

## 测试这些改进

```bash
# 运行验证器测试
go test ./internal/validator -v

# 运行基准测试
go test ./internal/validator -bench=.

# 测试速率限制
curl -X GET "http://localhost:8080/internal/search/keeps?q=test"

# 测试验证
curl -X GET "http://localhost:8080/internal/search/keeps?q=<script>alert('xss')</script>"

# 查看指标
curl http://localhost:8080/metrics
```

## 部署建议

1. **逐步部署**: 先在测试环境部署，验证所有功能正常
2. **监控指标**: 部署后密切监控新的指标
3. **调整限制**: 根据实际使用情况调整速率限制和验证规则
4. **性能测试**: 进行负载测试确保性能改进有效

这些改进将显著提升您应用的安全性、性能和可观测性。
