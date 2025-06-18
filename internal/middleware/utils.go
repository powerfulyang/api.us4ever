package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

// GetRealIP 优先获取 CF-Connecting-IP，其次 X-Forwarded-For，最后 c.IP()
func GetRealIP(c fiber.Ctx) string {
	if cfIP := c.Get("CF-Connecting-IP"); cfIP != "" {
		return cfIP
	}
	if xff := c.Get("X-Forwarded-For"); xff != "" {
		if idx := strings.Index(xff, ","); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}
	return c.IP()
}
