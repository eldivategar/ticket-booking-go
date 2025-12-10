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

	// Event routes
	eventGroup := v1.Group("/event")
	eventGroup.Use(deps.AuthMiddleware)
	eventGroup.Get("/:event_id", deps.EventHandler.GetEventByID)
	eventGroup.Get("/", deps.EventHandler.GetAllEvent)
	eventGroup.Post("/", deps.EventHandler.CreateEvent)
	eventGroup.Delete("/:event_id", deps.EventHandler.DeleteEvent)

	// Order routes
	orderGroup := v1.Group("/order")
	orderGroup.Use(deps.AuthMiddleware)
	orderGroup.Post("/", deps.OrderHandler.CreateOrder)
	orderGroup.Get("/:booking_id", deps.OrderHandler.GetOrderByBookingID)
	orderGroup.Get("/", deps.OrderHandler.GetOrderList)
	
	// Order Webhook routes
	orderWebhookGroup := v1.Group("/order/webhook")
	orderWebhookGroup.Post("/payment", deps.OrderHandler.ProcessPaymentWebhook)
}
