# Start with a Go base image
FROM golang:1.23.2-bullseye as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Install protoc and the Go plugin for protoc
RUN apt-get update && \
    apt-get install -y protobuf-compiler

RUN make generate

RUN go build -o myapp .

FROM debian:bullseye-slim

COPY --from=builder /app/myapp /myapp

EXPOSE 3000

CMD ["/myapp"]
