package app

import (
	"github.com/NonsoAmadi10/p2p-analysis/controllers"
	"github.com/NonsoAmadi10/p2p-analysis/services"
	"github.com/gofiber/fiber/v2"
)

func App() *fiber.App {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/node-info", controllers.GetMetrics)

	app.Get("/conn-metrics", controllers.GetConnMetrics)

	go services.ConnectionMetrics()
	return app
}
