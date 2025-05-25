package client

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

// NewAuthClient создает новый экземпляр клиента для Auth сервиса
func NewAuthClient(serverAddr string) (*AuthClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to gRPC server: %v", err)
		return nil, err
	}

	client := proto.NewAuthServiceClient(conn)

	return &AuthClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *AuthClient) Close() error {
	return c.conn.Close()
}

func (c *AuthClient) Register(ctx context.Context, username, email, password string) (*proto.AuthResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.client.Register(ctx, &proto.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	})
}

func (c *AuthClient) Login(ctx context.Context, email, password string) (*proto.AuthResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.client.Login(ctx, &proto.LoginRequest{
		Email:    email,
		Password: password,
	})
}

// GetUserInfo получает информацию о пользователе по ID
func (c *AuthClient) GetUserInfo(ctx context.Context, userID string) (*proto.UserInfoResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.client.GetUserInfo(ctx, &proto.UserInfoRequest{
		UserId: userID,
	})
}

func (c *AuthClient) UpdateProfile(ctx context.Context, userID, username, email string) (*proto.UpdateResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.client.UpdateProfile(ctx, &proto.UpdateProfileRequest{
		UserId:   userID,
		Username: username,
		Email:    email,
	})
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
