# Start with a Go base image
FROM golang:1.22.5 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Install protoc and the Go plugin for protoc
RUN apt-get update && \
    apt-get install -y protobuf-compiler && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

RUN make generate

RUN go build -o myapp .

FROM gcr.io/distroless/base-debian10

COPY --from=builder /app/myapp /myapp

EXPOSE 8080

CMD ["/myapp"]
