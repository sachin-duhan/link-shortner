package routes

import (
	"github.com/gofiber/fiber/v2"
)

func healthHandler(c *fiber.Ctx) error {
	return c.Status(200).SendString("Working")
}

func SetupRoutes(app *fiber.App) {
	app.Get("/health", healthHandler)
	app.Get("/:url", ResolveUrl)
	app.Post("/api/v1", ShortenUrl)
}