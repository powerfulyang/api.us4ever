package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"api.us4ever/internal/config"
	"api.us4ever/internal/logger"
	"api.us4ever/internal/metrics"
	"api.us4ever/internal/server"
	"api.us4ever/internal/task"
	"go.uber.org/zap"
)

const (
	// shutdownTimeout defines how long to wait for graceful shutdown
	shutdownTimeout = 5 * time.Second
	// defaultPort is used when no port is configured
	defaultPort = 8080
)

var (
	// Application loggers
	mainLogger      *logger.Logger
	shutdownLogger  *logger.Logger
	schedulerLogger *logger.Logger
	metricsLogger   *logger.Logger
)

func init() {
	var err error

	mainLogger, err = logger.New("main")
	if err != nil {
		panic("failed to initialize main logger: " + err.Error())
	}

	shutdownLogger, err = logger.New("shutdown")
	if err != nil {
		panic("failed to initialize shutdown logger: " + err.Error())
	}

	schedulerLogger, err = logger.New("scheduler")
	if err != nil {
		panic("failed to initialize scheduler logger: " + err.Error())
	}

	metricsLogger, err = logger.New("metrics")
	if err != nil {
		panic("failed to initialize metrics logger: " + err.Error())
	}
}

// gracefulShutdown handles the graceful shutdown of the server and scheduler
func gracefulShutdown(fiberServer *server.FiberServer, scheduler *task.Scheduler, done chan bool) {
	// Create context that listens for the interrupt signal from the OS
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal
	<-ctx.Done()

	shutdownLogger.Info("shutting down gracefully, press Ctrl+C again to force")

	// Create context with timeout for graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Shutdown the fiber server
	if err := fiberServer.ShutdownWithContext(shutdownCtx); err != nil {
		shutdownLogger.Error("server forced to shutdown with error",
			zap.Error(err),
		)
	}

	// Stop the task scheduler
	if scheduler != nil {
		scheduler.Stop()
	}

	shutdownLogger.Info("server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

// getPort returns the port to listen on, with fallback logic
func getPort(appConfig *config.AppConfig) int {
	if appConfig != nil && appConfig.Server.Port > 0 {
		return appConfig.Server.Port
	}

	if portStr := os.Getenv("PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil && port > 0 {
			return port
		}
	}

	return defaultPort
}

// getListenAddress returns the appropriate listen address based on environment
func getListenAddress(appConfig *config.AppConfig, port int) string {
	if appConfig != nil && logger.IsLocalDev(appConfig.AppEnv) {
		// Local development environment listens on localhost
		return fmt.Sprintf("localhost:%d", port)
	}
	// Other environments listen on 0.0.0.0, suitable for containers or servers
	return fmt.Sprintf("0.0.0.0:%d", port)
}

// initializeScheduler initializes and starts the task scheduler
func initializeScheduler(fiberServer *server.FiberServer) *task.Scheduler {
	scheduler, err := task.NewScheduler()
	if err != nil {
		schedulerLogger.Error("failed to initialize task scheduler",
			zap.Error(err),
		)
		return nil
	}

	// Register tasks
	if err := task.RegisterTasks(scheduler, fiberServer); err != nil {
		schedulerLogger.Error("failed to register tasks",
			zap.Error(err),
		)
		scheduler.Stop() // Clean up if registration fails
		return nil
	}

	// Start the scheduler
	scheduler.Start()
	schedulerLogger.Info("task scheduler started successfully")
	return scheduler
}

// initializeMetrics initializes the metrics collection
func initializeMetrics() {
	// 初始化指标收集器但不需要保存返回值
	_, err := metrics.StartMetricsCollection()
	if err != nil {
		metricsLogger.Error("failed to initialize metrics collector",
			zap.Error(err),
		)
		return
	}

	metricsLogger.Info("metrics collection started successfully")
}

func main() {
	// Initialize configuration
	appConfig := config.GetAppConfig()
	if appConfig == nil {
		mainLogger.Fatal("failed to load application configuration")
	}

	// Initialize fiber server
	fiberServer := server.New()
	if fiberServer == nil {
		mainLogger.Fatal("failed to initialize fiber server")
	}

	// Initialize metrics collection
	initializeMetrics()

	// Initialize task scheduler
	scheduler := initializeScheduler(fiberServer)

	// Register routes
	fiberServer.RegisterFiberRoutes()

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Start server in a goroutine
	go func() {
		port := getPort(appConfig)
		listenAddr := getListenAddress(appConfig, port)

		if err := fiberServer.Listen(listenAddr); err != nil {
			mainLogger.Fatal("failed to start server",
				zap.Error(err),
			)
		}
	}()

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(fiberServer, scheduler, done)

	// Wait for the graceful shutdown to complete
	<-done
	mainLogger.Info("graceful shutdown complete")
}
