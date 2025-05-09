package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"api.us4ever/internal/config"
	"api.us4ever/internal/server"
	"api.us4ever/internal/task"
)

func gracefulShutdown(fiberServer *server.FiberServer, scheduler *task.Scheduler, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := fiberServer.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	// 停止定时任务调度器
	scheduler.Stop()

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	// 初始化配置中心
	appConfig := config.GetAppConfig()

	fiberServer := server.New()
	// 初始化定时任务调度器
	scheduler, err := task.NewScheduler()
	if err != nil {
		log.Printf("初始化定时任务调度器失败: %v", err)
	} else {
		// 注册定时任务
		if err := task.RegisterTasks(scheduler, fiberServer); err != nil {
			log.Printf("注册定时任务失败: %v", err)
		} else {
			// 启动定时任务调度器
			scheduler.Start()
		}
	}
	fiberServer.RegisterFiberRoutes()

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	go func() {
		var port int
		if appConfig != nil {
			port = appConfig.Server.Port
		} else {
			port, _ = strconv.Atoi(os.Getenv("PORT"))
		}

		var listenAddr string
		if appConfig.AppEnv == "local" {
			// 本地开发环境监听 localhost
			listenAddr = fmt.Sprintf("localhost:%d", port)
		} else {
			// 其他环境监听 0.0.0.0，适合部署在容器或服务器上
			listenAddr = fmt.Sprintf("0.0.0.0:%d", port)
		}
		err := fiberServer.Listen(listenAddr)
		if err != nil {
			panic(fmt.Sprintf("http fiberServer error: %s", err))
		}
	}()

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(fiberServer, scheduler, done)

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}
