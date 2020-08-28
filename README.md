# fiber-http

A minimalist Golang web application for testing Opsani
optimizations with high throughput, predictable performance, and metrics instrumentation.

Built with [Fiber](https://docs.gofiber.io/) and [FastHTTP](https://github.com/valyala/fasthttp).

The app exposes several endpoints:

* `/` - Returns a 200 (Ok) `text/plain` respomse.
* `/metrics` - Metrics in Prometheus format for scraping.
* `/cpu{?duration}` - Consume CPU resources for the given duration (in Golang Duration string format). Default: `100ms`
* `/memory{?size}` - Consume memory resources by allocating a byte array of the given size (in human readable byte string format). Default: `10MB`
* `/time{?duration}` - Consume time by sleeping for the given duration (in Golang Duration string format). Default: `100ms`

The resource endpoints of `cpu`, `memory`, and `time` are useful for triggering the artificial consumption
of resources for testing autoscale behaviors, error handling, response to latency, etc.

## Instrumentation

The application includes support for instrumentation with metrics systems for testing optization
with different measurement backends.

A [Prometheus](https://prometheus.io/) metrics endpoint is exposed at `/metrics` that can be scraped.

[New Relic](https://newrelic.com/) instrumentation can be activated by setting the `NEW_ RELIC_LICENSE_KEY` environment
variable to a valid New Relic license key. The middleware will activate and log a status
message after initialization.

## Docker images

Docker images are published to [Docker Hub](https://hub.docker.com/r/opsani/fiber-http).

```console
$ docker pull opsani/fiber-http:latest
```

Tasks for working with the container image are in the
[Makefile](Makefile).

## GOMAXPROCS

This application utilizes [automaxprocs](https://github.com/uber-go/automaxprocs) to correctly
align the value of `GOMAXPROCS` when running in a container.

By default, `GOMAXPROCS` will align with the host core count rather than the CPU resources 
actually available to the container via the cgroups CPU quota. This misalignment will lead to 
resource exhaustion and make optimization impossible.

## License

Distributed under the terms of the Apache 2.0 license.
