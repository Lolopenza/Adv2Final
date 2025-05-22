package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"payment-service/internal/domain"
	"payment-service/internal/usecase"
	pb "payment-service/proto"
)

type subscriptionServer struct {
	pb.UnimplementedSubscriptionServiceServer
	subscriptionUseCase domain.SubscriptionUseCase
}

func NewSubscriptionServer(subscriptionUseCase domain.SubscriptionUseCase) pb.SubscriptionServiceServer {
	return &subscriptionServer{
		subscriptionUseCase: subscriptionUseCase,
	}
}

func (s *subscriptionServer) CreateSubscription(ctx context.Context, req *pb.CreateSubscriptionRequest) (*pb.CreateSubscriptionResponse, error) {
	subscription := &domain.Subscription{
		CustomerEmail: req.CustomerEmail,
		PlanName:      req.PlanName,
		Price:         req.Price,
		Currency:      req.Currency,
		StartDate:     time.Now(),
		EndDate:       time.Now().AddDate(0, 1, 0), // 1 month subscription
	}

	if err := s.subscriptionUseCase.CreateSubscription(subscription); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateSubscriptionResponse{
		SubscriptionId: subscription.ID,
		Status:         string(subscription.Status),
	}, nil
}

func (s *subscriptionServer) GetSubscription(ctx context.Context, req *pb.GetSubscriptionRequest) (*pb.GetSubscriptionResponse, error) {
	subscription, err := s.subscriptionUseCase.GetSubscription(req.SubscriptionId)
	if err != nil {
		if err == usecase.ErrSubscriptionNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetSubscriptionResponse{
		SubscriptionId: subscription.ID,
		CustomerEmail:  subscription.CustomerEmail,
		PlanName:       subscription.PlanName,
		Price:          subscription.Price,
		Currency:       subscription.Currency,
		Status:         string(subscription.Status),
		StartDate:      timestamppb.New(subscription.StartDate),
		EndDate:        timestamppb.New(subscription.EndDate),
		CreatedAt:      timestamppb.New(subscription.CreatedAt),
		UpdatedAt:      timestamppb.New(subscription.UpdatedAt),
	}, nil
}

func (s *subscriptionServer) CancelSubscription(ctx context.Context, req *pb.CancelSubscriptionRequest) (*pb.CancelSubscriptionResponse, error) {
	if err := s.subscriptionUseCase.CancelSubscription(req.SubscriptionId); err != nil {
		if err == usecase.ErrSubscriptionNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if err == usecase.ErrInvalidStatus {
			return nil, status.Error(codes.FailedPrecondition, "Only active subscriptions can be cancelled")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	subscription, err := s.subscriptionUseCase.GetSubscription(req.SubscriptionId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CancelSubscriptionResponse{
		Status: string(subscription.Status),
	}, nil
}

func (s *subscriptionServer) RenewSubscription(ctx context.Context, req *pb.RenewSubscriptionRequest) (*pb.RenewSubscriptionResponse, error) {
	if err := s.subscriptionUseCase.RenewSubscription(req.SubscriptionId); err != nil {
		if err == usecase.ErrSubscriptionNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if err == usecase.ErrInvalidStatus {
			return nil, status.Error(codes.FailedPrecondition, "Only active subscriptions can be renewed")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	subscription, err := s.subscriptionUseCase.GetSubscription(req.SubscriptionId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RenewSubscriptionResponse{
		Status: string(subscription.Status),
	}, nil
}

func (s *subscriptionServer) ListSubscriptions(ctx context.Context, req *pb.ListSubscriptionsRequest) (*pb.ListSubscriptionsResponse, error) {
	subscriptions, total, err := s.subscriptionUseCase.ListSubscriptions(req.CustomerEmail, req.Page, req.Limit)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &pb.ListSubscriptionsResponse{
		Total: total,
	}

	for _, subscription := range subscriptions {
		response.Subscriptions = append(response.Subscriptions, &pb.GetSubscriptionResponse{
			SubscriptionId: subscription.ID,
			CustomerEmail:  subscription.CustomerEmail,
			PlanName:       subscription.PlanName,
			Price:          subscription.Price,
			Currency:       subscription.Currency,
			Status:         string(subscription.Status),
			StartDate:      timestamppb.New(subscription.StartDate),
			EndDate:        timestamppb.New(subscription.EndDate),
			CreatedAt:      timestamppb.New(subscription.CreatedAt),
			UpdatedAt:      timestamppb.New(subscription.UpdatedAt),
		})
	}

	return response, nil
}
