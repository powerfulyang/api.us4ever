package routes

import (
	"github.com/gofiber/fiber/v2"
)

type SearchRoutes struct {
	app *fiber.App
}

func NewSearchRoutes(app *fiber.App) *SearchRoutes {
	return &SearchRoutes{
		app: app,
	}
}

func (r *SearchRoutes) Register() {
	internal := r.app.Group("/internal")
	searchGroup := internal.Group("/search")

	// 新的搜索路由
	searchGroup.Get("/keeps", r.searchKeepsHandler)
	searchGroup.Get("/moments", r.searchMomentsHandler)

	// 保持向后兼容的旧路由
	internal.Get("/keeps/search", r.searchKeepsHandler)
	internal.Get("/moments/search", r.searchMomentsHandler)
}

func (r *SearchRoutes) searchKeepsHandler(c *fiber.Ctx) error {
	// TODO: 实现搜索逻辑
	return c.SendString("Search keeps endpoint")
}

func (r *SearchRoutes) searchMomentsHandler(c *fiber.Ctx) error {
	// TODO: 实现搜索逻辑
	return c.SendString("Search moments endpoint")
}
