package routes

import (
	"github.com/gofiber/fiber/v3"
)

var routeRegistry = make([]func(*fiber.App), 0)

func Register(fn func(*fiber.App)) {
	routeRegistry = append(routeRegistry, fn)
}

func SetupAllRoutes(r *fiber.App) {
	for _, fn := range routeRegistry {
		fn(r)
	}
}
