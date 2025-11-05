package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupRoutes(app *fiber.App, deps *Dependencies) {
	// Define your application routes here
	app.Use(cors.New())
	
	// health check route
	app.Get("/ping", func (c *fiber.Ctx) error  {
		return c.SendString("Pong")
	})

	// API v1 group
	// v1 := app.Group("/api/v1")
}