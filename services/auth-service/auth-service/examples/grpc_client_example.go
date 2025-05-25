package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"auth-service/client"
	"auth-service/proto"
)

func main() {

	authClient, err := client.NewAuthClient("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to create Auth client: %v", err)
	}
	defer authClient.Close()

	ctx := context.Background()

	registerResp, err := authClient.Register(ctx, "test_user", "test@example.com", "password123")
	if err != nil {
		fmt.Printf("Registration error: %v\n", err)
	} else {
		printUserInfo("Registered user", registerResp.User, registerResp.Token)
	}

	loginResp, err := authClient.Login(ctx, "test@example.com", "password123")
	if err != nil {
		fmt.Printf("Login error: %v\n", err)
		os.Exit(1)
	}

	printUserInfo("Logged in user", loginResp.User, loginResp.Token)

	userInfoResp, err := authClient.GetUserInfo(ctx, loginResp.User.Id)
	if err != nil {
		fmt.Printf("Get user info error: %v\n", err)
	} else {
		printUserInfo("User info", userInfoResp.User, "")
	}

	updateResp, err := authClient.UpdateProfile(ctx, loginResp.User.Id, "updated_username", "")
	if err != nil {
		fmt.Printf("Update profile error: %v\n", err)
	} else {
		fmt.Printf("Update profile response: %s\n", updateResp.Message)
		printUserInfo("Updated user", updateResp.User, "")
	}

	changePassResp, err := authClient.ChangePassword(ctx, loginResp.User.Id, "password123", "newpassword123")
	if err != nil {
		fmt.Printf("Change password error: %v\n", err)
	} else {
		fmt.Printf("Change password response: %s\n", changePassResp.Message)
	}

	fmt.Println("gRPC client example completed!")
}

func printUserInfo(title string, user *proto.User, token string) {
	fmt.Printf("\n--- %s ---\n", title)
	fmt.Printf("ID: %s\n", user.Id)
	fmt.Printf("Username: %s\n", user.Username)
	fmt.Printf("Email: %s\n", user.Email)
	fmt.Printf("Created At: %s\n", user.CreatedAt.AsTime())
	fmt.Printf("Updated At: %s\n", user.UpdatedAt.AsTime())
	if token != "" {
		fmt.Printf("Token: %s\n", token)
	}
	fmt.Println("------------")
}
