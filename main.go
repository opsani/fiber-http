package main

import (
	"github.com/ansrivas/fiberprometheus"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
)

func main() {
	app := fiber.New()
	app.Use(middleware.Logger())
	app.Use(middleware.RequestID())

	prometheus := fiberprometheus.New("fiber-http")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	app.Get("/", func(c *fiber.Ctx) {
		c.Send("move along, nothing to see here")
	})

	app.Listen(8080)
}
