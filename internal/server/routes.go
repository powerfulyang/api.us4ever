package server

import (
	"time"

	"api.us4ever/internal/metrics"
	"api.us4ever/internal/middleware"
	"api.us4ever/internal/server/routes"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// 1. Recovery: 捕获后续所有中间件或处理器中的 panic，并将其转换为 500 错误。
	// 必须放在指标中间件之后，这样指标中间件才能捕获到它设置的 500 状态码。
	s.App.Use(middleware.NewRecoveryMiddleware())

	// 2. Metrics: 应该尽早放置，以测量整个请求的生命周期，包括其他中间件的执行时间。
	// 它在 c.Next() 调用之前启动计时器，在之后记录包括最终状态码在内的指标[3][4]。
	s.App.Use(metrics.NewMiddleware())

	// 3. Request ID: 为每个请求生成唯一ID，便于日志追踪。
	s.App.Use(requestid.New())

	// 4. Logging: 记录请求信息，可以利用前面生成的 Request ID。
	s.App.Use(middleware.NewLoggingMiddleware())

	// 5. Error Handler: 处理请求过程中产生的错误
	// 放在日志中间件之后，确保日志中间件能够记录到最终的状态码
	s.App.Use(middleware.NewErrorMiddleware())

	// 6. CORS: 处理跨域请求，通常放在业务逻辑之前。
	s.App.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Accept", "Authorization", "Content-Type"},
		MaxAge:       3600,
	}))

	// 7. Limiter: 在请求到达核心业务逻辑之前进行速率限制，保护应用。
	// 即使请求被限流（返回 429），指标中间件也因为在其之前注册而能记录到这次请求。
	s.App.Use(limiter.New(limiter.Config{
		Max:               30,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
		KeyGenerator: func(c fiber.Ctx) string {
			return middleware.GetRealIP(c)
		},
	}))

	// 注册基础路由
	routes.RegisterBaseRoutes(s.App)

	// 注册内部路由
	internalRoutes := routes.NewInternalRoutes(s.App, s.cfg, s.DbClient, s.EsClient)
	internalRoutes.Register()

	// 注册搜索路由
	searchRoutes := routes.NewSearchRoutes(s.App, s.EsClient, s.DbClient, s.KeepEsIndexAlias, s.MomentEsIndexAlias)
	searchRoutes.Register()

	// 注册重索引路由
	reindexRoutes := routes.NewReindexRoutes(s.App, s.EsClient, s.DbClient, s.KeepEsIndexAlias, s.MomentEsIndexAlias)
	reindexRoutes.Register()
}
