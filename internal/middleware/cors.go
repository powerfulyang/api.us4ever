package middleware

import "github.com/gofiber/fiber/v2"

// CORSMiddleware CORS Middleware
func CORSMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS,PATCH")
		c.Set("Access-Control-Allow-Headers", "Accept,Authorization,Content-Type")
		c.Set("Access-Control-Allow-Credentials", "false")
		c.Set("Access-Control-Max-Age", "300")

		return c.Next()
	}
}
