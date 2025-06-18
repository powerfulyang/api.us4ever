package routes

import (
	"fmt"
	"time"

	"api.us4ever/internal/logger"
	"api.us4ever/internal/metrics"
	"github.com/gofiber/fiber/v3"
)

var baseLogger, _ = logger.New("base")

func RegisterBaseRoutes(app *fiber.App) {
	app.Get("/", func(c fiber.Ctx) error {
		resp := fiber.Map{
			"message": "Hello World",
		}
		return c.JSON(resp)
	})
	app.Get("/panic", func(c fiber.Ctx) error {
		panic("error")
	})
	// 在默认设置下，Fiber 为了极致性能会复用内存。从 c.Params()、c.Body() 等方法获取的值，其底层字节缓冲区（buffer）会在请求结束后被下一个请求复用。
	// 使用命令 `curl http://localhost:8080/log/1 & curl http://localhost:8080/log/2 & curl http://localhost:8080/log/3` 测试
	app.Get("/log/:message", func(c fiber.Ctx) error {
		// 1. 从请求中获取参数
		msg := c.Params("message")
		// 正确用法
		// msg := utils.CloneString(c.Params("message"))
		baseLogger.Info(fmt.Sprintf("[Handler] Received message: %s", msg))

		// 2. 在一个新的 goroutine 中处理这个消息
		go func() {
			// 模拟耗时操作，如写入数据库
			time.Sleep(1 * time.Second)

			// 3. 在任务完成时打印消息
			// 这里会出现问题！
			baseLogger.Info(fmt.Sprintf("[Goroutine] Finished processing message: %s", msg))
		}()

		return c.SendString("Request accepted for: " + msg)
	})

	// 添加Prometheus指标端点
	app.Get("/metrics", metrics.GetMetricsHandler())
}
