package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/ansrivas/fiberprometheus"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/inhies/go-bytesize"

	_ "go.uber.org/automaxprocs"
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

	app.Get("/cpu", func(c *fiber.Ctx) {
		duration, err := time.ParseDuration(c.Query("duration", "100ms"))
		if err != nil {
			c.Next(err)
			return
		}

		consumeCPU(duration)
		c.Send(fmt.Sprintf("consumed CPU for %v\n", duration.String()))
	})

	app.Get("/memory", func(c *fiber.Ctx) {
		size, err := bytesize.Parse(c.Query("size", "10MB"))
		if err != nil {
			c.Next(err)
			return
		}

		data := make([]byte, size)
		c.Send(fmt.Sprintf("allocated %v (%d bytes) of memory\n", size.String(), len(data)))
	})

	app.Get("/time", func(c *fiber.Ctx) {
		duration, err := time.ParseDuration(c.Query("duration", "100ms"))
		if err != nil {
			c.Next(err)
			return
		}

		time.Sleep(duration)
		c.Send(fmt.Sprintf("slept for %v\n", duration.String()))
	})

	app.Listen(8080)
}

func consumeCPU(duration time.Duration) {
	maxProcs := runtime.GOMAXPROCS(0)
	stop := make(chan bool)

	for i := 0; i < maxProcs; i++ {
		go func() {
			for {
				select {
				case <-stop:
					return
				default:
				}
			}
		}()
	}

	time.Sleep(duration)

	for i := 0; i < maxProcs; i++ {
		stop <- true
	}
}
