package middleware

import (
	"sync"

	"github.com/gofiber/fiber/v2"
)

var mu sync.Mutex

func RaceCondition() fiber.Handler {
	return func(c *fiber.Ctx) error {
		mu.Lock()
		defer mu.Unlock()

		return c.Next()
	}
}