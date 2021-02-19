#FROM golang:1.14
#
#WORKDIR /go/src/app
#COPY . .
#
#RUN go get -d -v ./...
#RUN go install -v ./...
#
#CMD ["app"]

# Dockerfile from https://codefresh.io/docs/docs/learn-by-example/golang/golang-hello-world/
FROM golang:1.14 AS build_base

#RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /tmp/go-sample-app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Unit tests
#RUN CGO_ENABLED=0 go test -v

# Build the Go app
RUN CGO_ENABLED=0 go build -o ./out/go-sample-app .

# Start fresh from a smaller image
FROM alpine:3.9 
#RUN apk add ca-certificates

COPY --from=build_base /tmp/go-sample-app/out/go-sample-app /app/go-sample-app
COPY nodes.cfg .

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["/app/go-sample-app"]
