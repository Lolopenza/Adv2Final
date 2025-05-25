package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"auth-service/proto"
)

type AuthClient struct {
	client proto.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewAuthClient(serverAddr string) (*AuthClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to gRPC server: %v", err)
		return nil, err
	}
	return &AuthClient{
		client: proto.NewAuthServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *AuthClient) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) (*proto.ChangePasswordResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.client.ChangePassword(ctx, &proto.ChangePasswordRequest{
		UserId:          userID,
		CurrentPassword: currentPassword,
		NewPassword:     newPassword,
	})
}
