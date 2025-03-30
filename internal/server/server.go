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

	return server
}
