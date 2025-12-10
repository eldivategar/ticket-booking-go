package middleware

import (
	"go-war-ticket-service/internal/domain"
	"go-war-ticket-service/internal/platform/responses"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func RateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        20,
		Expiration: time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return responses.Error(c, fiber.StatusTooManyRequests, domain.ErrTooManyRequests.Error())
		},
	})
}
