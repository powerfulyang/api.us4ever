package routes

import (
	"github.com/gofiber/fiber/v2"
)

type ReindexRoutes struct {
	app *fiber.App
}

func NewReindexRoutes(app *fiber.App) *ReindexRoutes {
	return &ReindexRoutes{
		app: app,
	}
}

func (r *ReindexRoutes) Register() {
	internal := r.app.Group("/internal")

	// 重索引端点
	internal.Post("/keeps/reindex", r.reindexKeepsHandler)
	internal.Post("/moments/reindex", r.reindexMomentsHandler)
}

func (r *ReindexRoutes) reindexKeepsHandler(c *fiber.Ctx) error {
	// TODO: 实现重索引逻辑
	return c.SendString("Reindex keeps endpoint")
}

func (r *ReindexRoutes) reindexMomentsHandler(c *fiber.Ctx) error {
	// TODO: 实现重索引逻辑
	return c.SendString("Reindex moments endpoint")
}
