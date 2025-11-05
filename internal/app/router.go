package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupRoutes(app *fiber.App, deps *Dependencies) {
	// Define your application routes here
	app.Use(cors.New())

	// ping check route
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("Pong")
	})

	// API v1 group
	v1 := app.Group("/api/v1")

	// Auth routes
	authGroup := v1.Group("/auth")
	authGroup.Post("/register", deps.AuthHandler.Register)
	authGroup.Post("/login", deps.AuthHandler.Login)

	// User routes
	userGroup := v1.Group("/user")
	userGroup.Use(deps.AuthMiddleware)
	userGroup.Get("/me", deps.UserHandler.GetMyProfile)
}
