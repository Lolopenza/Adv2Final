package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"payment-service/cmd/config"
	"payment-service/internal/delivery/email"
	grpcDelivery "payment-service/internal/delivery/grpc"
	natsDelivery "payment-service/internal/delivery/nats"
	"payment-service/internal/repository/cache"
	"payment-service/internal/repository/postgres"
	"payment-service/internal/usecase"
	pb "payment-service/proto"

	"github.com/redis/go-redis/v9"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get configuration from environment
	port := config.GetEnv("GRPC_PORT", "50051")

	// Parse SMTP configuration
	smtpUser := parseEnvVar("SMTP_USER")
	smtpPass := parseEnvVar("SMTP_PASS")
	smtpHost := parseEnvVar("SMTP_HOST")
	smtpPort := parseEnvVar("SMTP_PORT")

	// Log SMTP config (mask password)
	maskedPass := "********"
	if len(smtpPass) > 0 {
		maskedPass = "********"
	}
	log.Printf("SMTP Configuration: User=%s, Pass=%s, Host=%s, Port=%s",
		smtpUser, maskedPass, smtpHost, smtpPort)

	// Initialize PostgreSQL connection
	db, err := initPostgres()
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Initialize Redis connection
	redisClient, err := initRedis()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize NATS connection
	natsConn, err := initNATS()
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}

	// Initialize repositories and services
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB from gorm.DB: %v", err)
	}
	paymentRepo := postgres.NewPaymentRepository(sqlDB)
	paymentCache := cache.NewRedisCache(redisClient)
	eventPublisher := natsDelivery.NewEventPublisher(natsConn)
	emailService := email.NewEmailService(
		smtpUser,
		smtpPass,
		smtpHost,
		smtpPort,
	)

	// Initialize use case
	paymentUseCase := usecase.NewPaymentUseCase(
		paymentRepo,
		paymentCache,
		eventPublisher,
		emailService,
	)

	// Initialize subscription repositories and services
	subscriptionRepo := postgres.NewSubscriptionRepository(sqlDB)
	subscriptionCache := cache.NewSubscriptionCache(redisClient)

	// Initialize subscription use case
	subscriptionUseCase := usecase.NewSubscriptionUseCase(
		subscriptionRepo,
		subscriptionCache,
		eventPublisher,
		emailService,
	)

	// Create gRPC server
	server := grpc.NewServer()

	// Register payment service
	pb.RegisterPaymentServiceServer(server, grpcDelivery.NewPaymentServer(paymentUseCase))

	// Register subscription service
	pb.RegisterSubscriptionServiceServer(server, grpcDelivery.NewSubscriptionServer(subscriptionUseCase))

	// Start listening
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Payment service gRPC server started on port %s", port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	log.Println("Shutting down server...")
	server.GracefulStop()
}

func initPostgres() (*gorm.DB, error) {
	dsn := os.Getenv("POSTGRES_DSN")
	db, err := gorm.Open(gormPostgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Run migrations
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Create payments table if it doesn't exist
	_, err = sqlDB.Exec(`
		CREATE TABLE IF NOT EXISTS payments (
			id VARCHAR(255) PRIMARY KEY,
			amount DECIMAL(10,2) NOT NULL,
			currency VARCHAR(3) NOT NULL,
			status VARCHAR(20) NOT NULL,
			customer_email VARCHAR(255) NOT NULL,
			description TEXT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_payments_customer_email ON payments(customer_email);
		CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);
		CREATE INDEX IF NOT EXISTS idx_payments_created_at ON payments(created_at);
		
		CREATE TABLE IF NOT EXISTS subscriptions (
			id VARCHAR(255) PRIMARY KEY,
			customer_email VARCHAR(255) NOT NULL,
			plan_name VARCHAR(255) NOT NULL,
			price DECIMAL(10,2) NOT NULL,
			currency VARCHAR(3) NOT NULL,
			status VARCHAR(20) NOT NULL,
			start_date TIMESTAMP NOT NULL,
			end_date TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_subscriptions_customer_email ON subscriptions(customer_email);
		CREATE INDEX IF NOT EXISTS idx_subscriptions_status ON subscriptions(status);
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func initRedis() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

func initNATS() (*nats.Conn, error) {
	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		return nil, err
	}

	return nc, nil
}

// parseEnvVar handles the multi-value environment variables
func parseEnvVar(key string) string {
	rawValue := os.Getenv(key)
	if rawValue == "" {
		return ""
	}

	// For SMTP_PASS specifically, remove all spaces (common with app passwords)
	if key == "SMTP_PASS" {
		// Remove all spaces from the password string
		return strings.ReplaceAll(rawValue, " ", "")
	}

	// If the value contains spaces, it might be multiple variables in one line
	if strings.Contains(rawValue, " ") {
		parts := strings.Fields(rawValue)
		for _, part := range parts {
			if strings.HasPrefix(part, key+"=") {
				return strings.TrimPrefix(part, key+"=")
			}
		}
		// If no explicit match, just return the first part
		return parts[0]
	}

	return rawValue
}
