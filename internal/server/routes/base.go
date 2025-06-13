package routes

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterBaseRoutes(app *fiber.App) {
	app.Get("/", HelloWorldHandler)
	app.Get("/err", ErrorHandler)
}

func HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}
	return c.JSON(resp)
}

func ErrorHandler(c *fiber.Ctx) error {
	panic("error")
}
