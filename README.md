# fiber-http

A minimalist Golang web application for testing Opsani optimizations with high
throughput, predictable performance, and metrics instrumentation.

Built with [Fiber](https://docs.gofiber.io/) and
[FastHTTP](https://github.com/valyala/fasthttp).

## Exposed endpoints

* `/` - Returns a 200 (Ok) `text/plain` response of "move along, nothing to see here".
* `/metrics` - Metrics in Prometheus format for scraping.
* `/cpu{?duration}` - Consume CPU resources for the given duration (in Golang
  Duration string format). Default: `100ms`
* `/memory{?size}` - Consume memory resources by allocating a byte array of the
  given size (in human readable byte string format). Default: `10MB`
* `/time{?duration}` - Consume time by sleeping for the given duration (in
  Golang Duration string format). Default: `100ms`
* `/request{?url}` - Proxy an HTTP GET request to a URL and return the status code & message body retrieved..

The resource endpoints of `cpu`, `memory`, and `time` are useful for triggering
the artificial consumption of resources for testing autoscale behaviors, error
handling, response to latency, etc.

The `request` endpoint enables testing of service dependencies and can be reentrantly
chained. For example, when running an instance locally on port 8480 a request made
to `http://localhost:8480/request?url=http://localhost:8480/time?duration=45ms` would
simulate an upstream service with a 45ms latency.

## Instrumentation

The application includes support for instrumentation with metrics systems for
testing optization with different measurement backends.

A [Prometheus](https://prometheus.io/) metrics endpoint is exposed at `/metrics`
that can be scraped.

[New Relic](https://newrelic.com/) instrumentation can be activated by setting
the `NEW_RELIC_LICENSE_KEY` environment variable to a valid New Relic license
key. The middleware will activate and log a status message after initialization.

Set `NEW_RELIC_APP_NAME` to define the corresponding New Relic APM identifier (default: `fiber-http`).

## Listening ports

By default, the server listens on HTTP port 8480 and HTTPS port 8843 (N + 8400). The port can be changed via the `HTTP_PORT` and `HTTPS_PORT` environment variables, respectively.

## Initial memory allocation

An initial resident memory allocation be created by setting the `INIT_MEMORY_SIZE` environment variable (in human readable byte string format).

## Docker images

Docker images are published to [Docker
Hub](https://hub.docker.com/r/opsani/fiber-http).

```console
$ docker pull opsani/fiber-http:latest
```

Tasks for working with the container image are in the [Makefile](Makefile).

## License

Distributed under the terms of the Apache 2.0 license.
