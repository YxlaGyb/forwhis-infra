package middleware

import "github.com/gofiber/fiber/v3"

func Cors() fiber.Handler {
	return func(c fiber.Ctx) error {
		origin := c.Get("Origin")
		if origin != "" {
			c.Set("Access-Control-Allow-Origin", origin)
		} else {
			c.Set("Access-Control-Allow-Origin", "*")
		}

		c.Set("Access-Control-Allow-Methods", "GET, POST")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Set("Access-Control-Expose-Headers", "X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset")
		c.Set("Access-Control-Max-Age", "600")
		c.Set("Access-Control-Allow-Credentials", "true")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(204)
		}

		return c.Next()
	}
}
