package main

import (
	"fmt"
	"log"
	"net"

	"appointment-service/configs"
	grpcHandler "appointment-service/internal/delivery/grpc"
	httpHandler "appointment-service/internal/delivery/http"
	"appointment-service/internal/repository/postgres"
	"appointment-service/internal/usecase"
	pb "appointment-service/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg := configs.LoadConfig()

	// Initialize repository
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.Password, cfg.Database.DBName, cfg.Database.SSLMode)
	repo, err := postgres.NewPostgresRepository(dsn)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	// Initialize use case
	useCase, err := usecase.NewAppointmentUseCase(repo, cfg.NATS.URL, fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port))
	if err != nil {
		log.Fatalf("Failed to initialize use case: %v", err)
	}

	// Initialize HTTP server
	router := gin.Default()
	handler := httpHandler.NewHandler(useCase)
	handler.RegisterRoutes(router)

	// Start HTTP server
	go func() {
		if err := router.Run(":" + cfg.Server.HTTPPort); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Initialize gRPC server
	lis, err := net.Listen("tcp", ":"+cfg.Server.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	grpcHandler := grpcHandler.NewServer(useCase)
	pb.RegisterAppointmentServiceServer(grpcServer, grpcHandler)

	// Start gRPC server
	log.Printf("Starting gRPC server on port %s", cfg.Server.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
