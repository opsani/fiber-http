package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/ansrivas/fiberprometheus"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/inhies/go-bytesize"

	"github.com/newrelic/go-agent/v3/newrelic"
)

func main() {
	app := fiber.New()
	app.Use(middleware.Logger())
	app.Use(middleware.RequestID())

	prometheus := fiberprometheus.New("fiber-http")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

  listenPort := 8080
	listenPortInput := os.Getenv("FIBER_HTTP_LISTEN_PORT")
	if listenPortInput != "" {
		lp, err := strconv.Atoi(listenPortInput)
		if err != nil {
			log.Printf("ERROR setting FIBER_HTTP_LISTEN_PORT: %s\n", err)
		} else {
			listenPort = lp
		}
	}

  newrelicAppName := os.Getenv("NEW_RELIC_APP_NAME")
	if newrelicAppName == "" {
		newrelicAppName = "fiber-http"
	}
	// activate New Relic if NEW_RELIC_LICENSE_KEY is in the environment
	if newrelicLicenseKey := os.Getenv("NEW_RELIC_LICENSE_KEY"); newrelicLicenseKey != "" {
		newrelicApp, err := newrelic.NewApplication(
			newrelic.ConfigAppName(newrelicAppName),
			newrelic.ConfigLicense(newrelicLicenseKey),
		)
		if err == nil {
			app.Use(NewRelicMiddleware(newrelicApp))
			log.Println("New Relic middleware initialized")
		} else {
			log.Printf("WARNING: failed to initialize New Relic: %s\n", err)
		}
	}

	app.Get("/", func(c *fiber.Ctx) {
		c.Send("move along, nothing to see here")
	})

  app.Get("/call", func(c *fiber.Ctx) {
		remoteURL := c.Query("url")
		if remoteURL == "" {
			c.Send("no url read, nothing to see here")
			return
		}
		out, err := exec.Command("curl", remoteURL).Output()
    if err != nil {
			c.Send(err)
		}
		c.Send(out)
  })

	app.Get("/cpu", func(c *fiber.Ctx) {
		duration, err := time.ParseDuration(c.Query("duration", "100ms"))
		if err != nil {
			c.Next(err)
			return
		}

		x := 0.0001
		start := time.Now()
		for time.Since(start) < duration {
			x += math.Sqrt(x)
		}

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

	app.Listen(listenPort)
}

// NewRelicMiddleware instruments the request with New Relic
func NewRelicMiddleware(app *newrelic.Application) fiber.Handler {
	return func(c *fiber.Ctx) {
		// start an HTTP transaction with New Relic
		txn := app.StartTransaction(c.Path())
		defer txn.End()

		// let Fiber process the request
		c.Next()

		// translate the FastHTTP request & response for New Relic
		hdr := make(http.Header)
		c.Fasthttp.Request.Header.VisitAll(func(k, v []byte) {
			sk := string(k)
			sv := string(v)
			hdr.Set(sk, sv)
		})

		txn.SetWebRequest(newrelic.WebRequest{
			Header:    http.Header{},
			URL:       &url.URL{Path: c.Path()},
			Method:    c.Method(),
			Transport: newrelic.TransportHTTP,
		})

		// Get a New Relic wrapper for the response writer
		rw := txn.SetWebResponse(nil)
		rw.WriteHeader(c.Fasthttp.Response.StatusCode())
		rw.Write(c.Fasthttp.Response.Body())
	}
}
