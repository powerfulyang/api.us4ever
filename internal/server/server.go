package server

import (
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

	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "api.us4ever",
			AppName:      "api.us4ever",
		}),

		db:  database.New(),
		cfg: appConfig,
	}

	return server
}
