package server

import (
	"log"

	"api.us4ever/internal/config"
	"api.us4ever/internal/database"
	"github.com/gofiber/fiber/v2"
)

type FiberServer struct {
	*fiber.App

	db  database.Service
	cfg *config.AppConfig
}

func New() *FiberServer {
	// 获取配置
	appConfig := config.GetAppConfig()

	// 初始化数据库服务
	db, err := database.New()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "api.us4ever",
			AppName:      "api.us4ever",
		}),

		db:  db,
		cfg: appConfig,
	}

	// 注册配置变更回调，当配置变更时刷新数据库连接
	config.RegisterChangeCallback(server.handleConfigChange)

	return server
}

// handleConfigChange 处理配置变更
func (s *FiberServer) handleConfigChange(newConfig *config.AppConfig) {
	log.Println("配置变更，更新服务...")

	// 更新服务器配置
	s.cfg = newConfig

	// 刷新数据库连接
	if err := s.RefreshDatabase(); err != nil {
		log.Printf("更新数据库连接失败: %v", err)
	} else {
		log.Println("数据库连接已更新")
	}
}

// RefreshDatabase 重新创建数据库连接
func (s *FiberServer) RefreshDatabase() error {
	// 不使用类型断言，直接创建新连接
	newDb, err := database.New()
	if err != nil {
		return err
	}

	// 如果有旧连接，尝试关闭
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			log.Printf("Warning: error closing previous database connection: %v", err)
		}
	}

	// 更新连接
	s.db = newDb
	return nil
}
