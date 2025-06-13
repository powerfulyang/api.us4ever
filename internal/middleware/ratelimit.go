package middleware

import (
	"fmt"
	"sync"
	"time"

	"api.us4ever/internal/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// RateLimitConfig holds the configuration for rate limiting
type RateLimitConfig struct {
	// RequestsPerSecond defines the rate limit
	RequestsPerSecond int

	// BurstSize defines the burst capacity
	BurstSize int

	// KeyGenerator generates the key for rate limiting (e.g., by IP, user ID)
	KeyGenerator func(*fiber.Ctx) string

	// OnLimitReached is called when rate limit is exceeded
	OnLimitReached func(*fiber.Ctx) error

	// Logger for rate limit events
	Logger *logger.Logger
}

// RateLimiter manages rate limiting for different keys
type RateLimiter struct {
	config   RateLimitConfig
	limiters sync.Map // map[string]*rate.Limiter
	mu       sync.RWMutex
	logger   *logger.Logger
}

// NewRateLimitMiddleware creates a new rate limiting middleware
func NewRateLimitMiddleware(config RateLimitConfig) fiber.Handler {
	// Set defaults
	if config.RequestsPerSecond <= 0 {
		config.RequestsPerSecond = 10 // Default: 10 requests per second
	}

	if config.BurstSize <= 0 {
		config.BurstSize = config.RequestsPerSecond * 2 // Default: 2x the rate
	}

	if config.KeyGenerator == nil {
		config.KeyGenerator = defaultKeyGenerator
	}

	if config.OnLimitReached == nil {
		config.OnLimitReached = defaultLimitReachedHandler
	}

	if config.Logger == nil {
		var err error
		config.Logger, err = logger.New("ratelimit")
		if err != nil {
			panic("failed to create rate limit logger: " + err.Error())
		}
	}

	limiter := &RateLimiter{
		config: config,
		logger: config.Logger,
	}

	// Start cleanup goroutine
	go limiter.cleanup()

	return limiter.Handler
}

// Handler is the fiber middleware handler
func (rl *RateLimiter) Handler(c *fiber.Ctx) error {
	key := rl.config.KeyGenerator(c)

	// Get or create limiter for this key
	limiterInterface, _ := rl.limiters.LoadOrStore(key, rate.NewLimiter(
		rate.Limit(rl.config.RequestsPerSecond),
		rl.config.BurstSize,
	))

	limiter := limiterInterface.(*rate.Limiter)

	// Check if request is allowed
	if !limiter.Allow() {
		rl.logger.Warn("rate limit exceeded",
			zap.String("key", key),
			zap.String("ip", c.IP()),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
		)

		return rl.config.OnLimitReached(c)
	}

	// Log successful requests (debug level)
	rl.logger.Debug("request allowed",
		zap.String("key", key),
		zap.String("ip", c.IP()),
		zap.String("method", c.Method()),
		zap.String("path", c.Path()),
	)

	return c.Next()
}

// cleanup removes old limiters to prevent memory leaks
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.limiters.Range(func(key, value interface{}) bool {
			limiter := value.(*rate.Limiter)

			// Remove limiter if it hasn't been used recently
			// This is a simple heuristic - in production you might want more sophisticated cleanup
			if limiter.Tokens() == float64(rl.config.BurstSize) {
				rl.limiters.Delete(key)
				rl.logger.Debug("cleaned up unused rate limiter",
					zap.Any("key", key),
				)
			}

			return true
		})
	}
}

// defaultKeyGenerator generates a key based on client IP
func defaultKeyGenerator(c *fiber.Ctx) string {
	return c.IP()
}

// defaultLimitReachedHandler returns a 429 Too Many Requests response
func defaultLimitReachedHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
		"error": fiber.Map{
			"type":    "RateLimitError",
			"message": "Too many requests, please try again later",
			"code":    429,
		},
		"retry_after": "60s",
	})
}

// NewIPRateLimiter creates a rate limiter based on IP address
func NewIPRateLimiter(requestsPerSecond int) fiber.Handler {
	return NewRateLimitMiddleware(RateLimitConfig{
		RequestsPerSecond: requestsPerSecond,
		KeyGenerator:      defaultKeyGenerator,
	})
}

// NewUserRateLimiter creates a rate limiter based on user ID
func NewUserRateLimiter(requestsPerSecond int) fiber.Handler {
	return NewRateLimitMiddleware(RateLimitConfig{
		RequestsPerSecond: requestsPerSecond,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Try to get user ID from context or headers
			if userID := c.Locals("user_id"); userID != nil {
				return fmt.Sprintf("user:%v", userID)
			}

			// Fallback to IP if no user ID
			return fmt.Sprintf("ip:%s", c.IP())
		},
	})
}

// NewPathBasedRateLimiter creates different rate limits for different paths
func NewPathBasedRateLimiter(requestsPerSecond int) fiber.Handler {
	return NewRateLimitMiddleware(RateLimitConfig{
		RequestsPerSecond: requestsPerSecond,
		KeyGenerator: func(c *fiber.Ctx) string {
			path := c.Path()
			method := c.Method()
			return fmt.Sprintf("ip:%s,method:%s,path:%s", c.IP(), method, path)
		},
	})
}
