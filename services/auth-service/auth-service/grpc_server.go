package main

import (
	"context"
	"log"
	"net"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"auth-service/models"
	"auth-service/proto"
)

type grpcServer struct {
	proto.UnimplementedAuthServiceServer
}

func userToProto(user models.User) *proto.User {
	return &proto.User{
		Id:        user.ID.Hex(),
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func (s *grpcServer) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.AuthResponse, error) {

	var existingUser models.User
	err := userCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "User with this email already exists")
	}

	hashedPassword := hashPassword(req.Password)

	now := time.Now()
	newUser := models.User{
		ID:        primitive.NewObjectID(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err = userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to register user")
	}

	token, err := generateToken(newUser.ID.Hex())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate token")
	}

	return &proto.AuthResponse{
		Token: token,
		User:  userToProto(newUser),
	}, nil
}

func (s *grpcServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.AuthResponse, error) {

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Invalid credentials")
	}

	if user.Password != hashPassword(req.Password) {
		return nil, status.Errorf(codes.PermissionDenied, "Invalid credentials")
	}

	token, err := generateToken(user.ID.Hex())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate token")
	}

	return &proto.AuthResponse{
		Token: token,
		User:  userToProto(user),
	}, nil
}

func (s *grpcServer) GetUserInfo(ctx context.Context, req *proto.UserInfoRequest) (*proto.UserInfoResponse, error) {
	objID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid user ID")
	}

	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	return &proto.UserInfoResponse{
		User: userToProto(user),
	}, nil
}

func (s *grpcServer) UpdateProfile(ctx context.Context, req *proto.UpdateProfileRequest) (*proto.UpdateResponse, error) {
	objID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid user ID")
	}

	if req.Email != "" {
		var existingUser models.User
		err := userCollection.FindOne(ctx, bson.M{
			"email": req.Email,
			"_id":   bson.M{"$ne": objID},
		}).Decode(&existingUser)
		if err == nil {
			return nil, status.Errorf(codes.AlreadyExists, "Email is already taken")
		}
	}

	updateFields := bson.M{}
	if req.Username != "" {
		updateFields["username"] = req.Username
	}
	if req.Email != "" {
		updateFields["email"] = req.Email
	}
	updateFields["updated_at"] = time.Now()

	_, err = userCollection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": updateFields},
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update profile")
	}

	var updatedUser models.User
	err = userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&updatedUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to fetch updated user")
	}

	return &proto.UpdateResponse{
		Message: "Profile updated successfully",
		User:    userToProto(updatedUser),
	}, nil
}

func (s *grpcServer) ChangePassword(ctx context.Context, req *proto.ChangePasswordRequest) (*proto.ChangePasswordResponse, error) {
	objID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid user ID")
	}

	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	// Проверка текущего пароля
	if hashPassword(req.CurrentPassword) != user.Password {
		return nil, status.Errorf(codes.PermissionDenied, "Current password is incorrect")
	}

	// Обновление пароля
	hashedNewPassword := hashPassword(req.NewPassword)
	_, err = userCollection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{
			"password":   hashedNewPassword,
			"updated_at": time.Now(),
		}},
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update password")
	}

	return &proto.ChangePasswordResponse{
		Message: "Password updated successfully",
	}, nil
}

func startGRPCServer() error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	proto.RegisterAuthServiceServer(s, &grpcServer{})

	reflection.Register(s)

	log.Println("gRPC Server running on :50051")
	return s.Serve(lis)
}
