# fiber-http

A minimalist Golang web application for testing Opsani
optimizations with minimal dependencies, high throughput,
and predictable performance.

Built with [Fiber](https://docs.gofiber.io/) and [FastHTTP](https://github.com/valyala/fasthttp).

The app exposes two endpoints:

* `/` - Returns a 200 (Ok) `text/plain` respomse.
* `/metrics` - Metrics in Prometheus format for scraping.

## Docker images

Docker images are published to [Docker Hub](https://hub.docker.com/r/opsani/fiber-http).

Tasks for working with the container image are in the
[Makefile](Makefile).

## License

Distributed under the terms of the Apache 2.0 license.
