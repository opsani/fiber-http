# fiber-http

A minimalist Golang web application for testing Opsani optimizations with high
throughput, predictable performance, and metrics instrumentation.

Built with [Fiber](https://docs.gofiber.io/) and
[FastHTTP](https://github.com/valyala/fasthttp).

## Listening port
Set environment variable FIBER_HTTP_LISTEN_PORT to the desired TCP listening port.

The default listening port is 8080.

## Exposed endpoints

* `/` - Returns a 200 (Ok) `text/plain` respomse.
* `/metrics` - Metrics in Prometheus format for scraping.
* `/call{?url}` - Have fiber-http to call out to another URL.
* `/cpu{?duration}` - Consume CPU resources for the given duration (in Golang
  Duration string format). Default: `100ms`
* `/memory{?size}` - Consume memory resources by allocating a byte array of the
  given size (in human readable byte string format). Default: `10MB`
* `/time{?duration}` - Consume time by sleeping for the given duration (in
  Golang Duration string format). Default: `100ms`

## Remote example call to fiber-http server
```% curl "http://localhost:8080/call?url=http://localhost:8000/cpu?duration=250ms"
consumed CPU for 250ms```

The resource endpoints of `cpu`, `memory`, and `time` are useful for triggering
the artificial consumption of resources for testing autoscale behaviors, error
handling, response to latency, etc.

## Instrumentation

The application includes support for instrumentation with metrics systems for
testing optization with different measurement backends.

A [Prometheus](https://prometheus.io/) metrics endpoint is exposed at `/metrics`
that can be scraped.

[New Relic](https://newrelic.com/) instrumentation can be activated by setting
the `NEW_RELIC_LICENSE_KEY` environment variable to a valid New Relic license
key. The middleware will activate and log a status message after initialization.

Set `NEW_RELIC_APP_NAME` to define the corresponding New Relic APM identifier (default: `fiber-http`).

## Docker images

Docker images are published to [Docker
Hub](https://hub.docker.com/r/opsani/fiber-http).

```console
$ docker pull opsani/fiber-http:latest
```

Tasks for working with the container image are in the [Makefile](Makefile).

## License

Distributed under the terms of the Apache 2.0 license.
