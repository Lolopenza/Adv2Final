# Payment Service

A microservice for handling payment operations in a distributed system. Built with Clean Architecture principles and following best practices for microservices.

## Features

- gRPC API for payment operations
- PostgreSQL for persistent storage
- Redis for caching
- NATS for event-driven communication
- SMTP email notifications
- Clean Architecture implementation

## Prerequisites

- Go 1.21 or later
- PostgreSQL
- Redis
- NATS
- SMTP server (Gmail or Microsoft)

## Setup

1. Clone the repository:
```bash
git clone https://github.com/yourusername/payment-service.git
cd payment-service
```

2. Install dependencies:
```bash
go mod download
```

3. Copy the environment file and update the values:
```bash
cp .env.example .env
```

4. Update the `.env` file with your configuration values.

5. Run database migrations:
```bash
# TODO: Add migration commands
```

## Running the Service

```bash
go run cmd/main.go
```

The service will start on port 50051 by default.

## API Endpoints

The service exposes the following gRPC endpoints:

1. `InitiatePayment`: Start a new payment process
2. `ConfirmPayment`: Confirm a pending payment
3. `GetPaymentStatus`: Get the current status of a payment
4. `RefundPayment`: Process a refund for a completed payment

## Event Topics

The service publishes the following events to NATS:

- `payment.confirmed`: When a payment is confirmed
- `payment.failed`: When a payment fails
- `payment.refunded`: When a payment is refunded

## Testing

Run the test suite:

```bash
go test ./...
```

## Project Structure

```
.
├── cmd/
│   └── main.go
├── internal/
│   ├── domain/
│   │   └── payment.go
│   ├── usecase/
│   │   └── payment_usecase.go
│   ├── repository/
│   │   ├── postgres/
│   │   │   ├── migrations/
│   │   │   └── payment_repository.go
│   │   └── cache/
│   │       └── redis_cache.go
│   └── delivery/
│       ├── grpc/
│       │   └── payment_server.go
│       ├── nats/
│       │   └── event_publisher.go
│       └── email/
│           └── email_service.go
├── proto/
│   └── payment.proto
├── go.mod
├── go.sum
└── README.md
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 