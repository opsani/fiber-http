// Copyright 2020 Opsani
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/inhies/go-bytesize"
	"github.com/valyala/fasthttp"

	"github.com/newrelic/go-agent/v3/newrelic"
)

var once sync.Once
var app *fiber.App
var initMemory []byte

func newApp() *fiber.App {
	once.Do(func() {
		app = fiber.New()
		app.Use(logger.New())
		app.Use(requestid.New())

		prometheus := fiberprometheus.New("fiber-http")
		prometheus.RegisterAt(app, "/metrics")
		app.Use(prometheus.Middleware)

		// activate New Relic if NEW_RELIC_LICENSE_KEY is in the environment
		if newrelicLicenseKey := os.Getenv("NEW_RELIC_LICENSE_KEY"); newrelicLicenseKey != "" {
			newrelicAppName := os.Getenv("NEW_RELIC_APP_NAME")
			if newrelicAppName == "" {
				newrelicAppName = "fiber-http"
			}
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

		app.Get("/", func(c *fiber.Ctx) error {
			return c.SendString("move along, nothing to see here")
		})

		app.Get("/cpu", func(c *fiber.Ctx) error {
			operations, err := strconv.ParseUint(c.Query("operations", "0"), 10, 64)
			if err != nil {
				return err
			}
			duration, err := time.ParseDuration(c.Query("duration", "100ms"))
			if err != nil {
				return err
			}

			i := uint64(0)
			x := 0.0001
			start := time.Now()
			for time.Since(start) < duration {
				if operations != 0 && i == operations {
					break
				}
				x += math.Sqrt(x)
				i++
			}

			runtime := time.Since(start)
			return c.SendString(fmt.Sprintf("consumed CPU for %v operations in %v\n", i, runtime.String()))
		})

		app.Get("/memory", func(c *fiber.Ctx) error {
			size, err := bytesize.Parse(c.Query("size", "10MB"))
			if err != nil {
				return err
			}

			data := append([]byte{}, make([]byte, size)...)
			return c.SendString(fmt.Sprintf("allocated %v (%d bytes) of memory\n", size.String(), len(data)))
		})

		app.Get("/time", func(c *fiber.Ctx) error {
			duration, err := time.ParseDuration(c.Query("duration", "100ms"))
			if err != nil {
				return err
			}

			time.Sleep(duration)
			return c.SendString(fmt.Sprintf("slept for %v\n", duration.String()))
		})

		app.Get("/request", func(c *fiber.Ctx) error {
			remoteURL := c.Query("url")
			if remoteURL == "" {
				c.Status(fiber.StatusBadRequest)
				return c.SendString("error: missing required query parameter \"url\"")
			}

			client := fasthttp.Client{}
			statusCode, body, err := client.Get(nil, remoteURL)
			c.Status(statusCode)
			if err != nil {
				return err
			}
			return c.Send(body)
		})

		app.Use(func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusNotFound)
		})
	})

	return app
}

func main() {
	// Allocate an initial heap if requested
	if sizeEnv := os.Getenv("INIT_MEMORY_SIZE"); sizeEnv != "" {
		size, err := bytesize.Parse(sizeEnv)
		if err != nil {
			log.Fatal(err)
		}
		initMemory = append(initMemory, make([]byte, size)...)
		log.Printf("NOTICE: allocated %v (%d bytes) of memory\n", size.String(), len(initMemory))
	}

	httpPort := ":8480"
	if p := os.Getenv("HTTP_PORT"); p != "" {
		httpPort = p
	}
	app := newApp()

	// Load TLS assets
	cer, err := tls.LoadX509KeyPair("certs/dev.opsani.com+3.pem", "certs/dev.opsani.com+3-key.pem")
	if err != nil {
		log.Fatal(err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}

	// Create TLS port listener
	httpsPort := ":8843"
	if p := os.Getenv("HTTPS_PORT"); p != "" {
		httpsPort = p
	}
	ln, err := tls.Listen("tcp", httpsPort, config)
	if err != nil {
		panic(err)
	}

	// Listen with TLS on HTTPS_PORT (:8843)
	go func() {
		log.Fatal(app.Listener(ln))
	}()

	// Listen on HTTP_PORT (:8480)
	log.Fatal(app.Listen(httpPort))
}

// NewRelicMiddleware instruments the request with New Relic
func NewRelicMiddleware(app *newrelic.Application) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// start an HTTP transaction with New Relic
		txn := app.StartTransaction(c.Path())
		defer txn.End()

		// let Fiber process the request
		c.Next()

		// translate the FastHTTP request & response for New Relic
		hdr := make(http.Header)
		c.Context().Request.Header.VisitAll(func(k, v []byte) {
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
		rw.WriteHeader(c.Context().Response.StatusCode())
		_, err := rw.Write(c.Context().Response.Body())
		return err
	}
}
