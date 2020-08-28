# fiber-http

A minimalist Golang web application for testing Opsani
optimizations with minimal dependencies, high throughput,
and predictable performance.

Built with [Fiber](https://docs.gofiber.io/) and [FastHTTP](https://github.com/valyala/fasthttp).

The app exposes several endpoints:

* `/` - Returns a 200 (Ok) `text/plain` respomse.
* `/metrics` - Metrics in Prometheus format for scraping.
* `/cpu{?duration}` - Consume CPU resources for the given duration (in Golang Duration string format). Default: `100ms`
* `/memory{?size}` - Consume memory resources by allocating a byte array of the given size (in human readable byte string format). Default: `10MB`
* `/time{?duration}` - Consume time by sleeping for the given duration (in Golang Duration string format). Default: `100ms`

The resource endpoints of `cpu`, `memory`, and `time` are useful for triggering the artificial consumption
of resources for testing autoscale behaviors, error handling, response to latency, etc.

## Docker images

Docker images are published to [Docker Hub](https://hub.docker.com/r/opsani/fiber-http).

```console
$ docker pull opsani/fiber-http:latest
```

Tasks for working with the container image are in the
[Makefile](Makefile).

## License

Distributed under the terms of the Apache 2.0 license.
