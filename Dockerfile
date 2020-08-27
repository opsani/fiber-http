FROM golang:1.15-alpine AS build

WORKDIR /src/
COPY main.go go.* /src/
RUN go get -v -d .
RUN go build -o /bin/fiber-http

FROM alpine
COPY --from=build /bin/fiber-http /bin/fiber-http
ENTRYPOINT ["/bin/fiber-http"]
