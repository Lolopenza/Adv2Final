# Appointment Service

A microservice for managing appointment slots and bookings, built with Go using Clean Architecture.

## Features

- Create and manage time slots for businesses
- Book available slots for users
- REST API and gRPC interfaces
- Redis caching for improved performance
- NATS messaging for event-driven architecture
- PostgreSQL database for data persistence

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- PostgreSQL
- Redis
- NATS
- Protocol Buffers compiler (protoc)

## Configuration

The service can be configured using environment variables:

```bash
# Server
HTTP_PORT=8080
GRPC_PORT=9090

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=appointments
DB_SSLMODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# NATS
NATS_URL=nats://localhost:4222
```

## Running the Service

### Using Docker

1. Build and start the service with Docker Compose:

```bash
docker-compose up --build
```

### Running Locally

1. Install dependencies:

```bash
go mod download
```

2. Generate protobuf code:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/appointment.proto
```

3. Run the service:

```bash
go run cmd/main.go
```

## API Documentation

### REST API

#### Create Slots
```http
POST /api/v1/slots
Content-Type: application/json

{
  "business_id": "uuid",
  "date": "2025-06-01",
  "slots": ["09:00", "09:30", "10:00"]
}
```

#### Get Available Slots
```http
GET /api/v1/slots?business_id=uuid&date=2025-06-01
```

#### Book Appointment
```http
POST /api/v1/appointments
Content-Type: application/json

{
  "user_id": "uuid",
  "slot_id": "uuid"
}
```

### gRPC

The service exposes the following gRPC methods:

- `CreateSlot`: Create new time slots
- `ListSlots`: Get available slots
- `BookAppointment`: Book an appointment

## Testing

### Testing REST API

```bash
# Create slots
curl -X POST http://localhost:8080/api/v1/slots \
  -H "Content-Type: application/json" \
  -d '{"business_id":"uuid","date":"2025-06-01","slots":["09:00","09:30"]}'

# Get available slots
curl "http://localhost:8080/api/v1/slots?business_id=uuid&date=2025-06-01"

# Book appointment
curl -X POST http://localhost:8080/api/v1/appointments \
  -H "Content-Type: application/json" \
  -d '{"user_id":"uuid","slot_id":"uuid"}'
```

### Testing gRPC

Use a gRPC client like [grpcurl](https://github.com/fullstorydev/grpcurl) to test the gRPC endpoints:

```bash
# List available slots
grpcurl -plaintext -d '{"business_id":"uuid","date":"2025-06-01"}' \
  localhost:9090 appointment.AppointmentService/ListSlots
```

### Checking NATS Messages

To monitor NATS messages:

```bash
# Install NATS CLI
go install github.com/nats-io/natscli/nats@latest

# Subscribe to appointment events
nats sub appointment.created
```

### Checking Redis Cache

To inspect Redis cache:

```bash
# Connect to Redis CLI
redis-cli

# Check cached slots
GET slots:uuid:2025-06-01
```

## Project Structure

```
appointment-service/
├── cmd/
│   └── main.go
├── internal/
│   ├── delivery/
│   │   ├── grpc/
│   │   └── http/
│   ├── domain/
│   ├── repository/
│   ├── usecase/
│   └── model/
├── migrations/
├── proto/
├── configs/
├── Dockerfile
└── go.mod
``` 