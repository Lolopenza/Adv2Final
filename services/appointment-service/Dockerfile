# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install protoc and required tools
RUN apk add --no-cache protobuf-dev make git

# Copy go mod and sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Install protoc plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate protobuf code
RUN protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/appointment.proto

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o appointment-service ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/appointment-service .

# Expose ports
EXPOSE 8080 9090

# Run the application
CMD ["./appointment-service"]
