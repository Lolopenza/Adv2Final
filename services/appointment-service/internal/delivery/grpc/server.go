package grpc

import (
	"context"
	"time"

	"appointment-service/internal/domain"
	"appointment-service/internal/usecase"
	pb "appointment-service/proto"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedAppointmentServiceServer
	useCase usecase.UseCase
}

func NewServer(useCase usecase.UseCase) *Server {
	return &Server{
		useCase: useCase,
	}
}

func (s *Server) CreateSlot(ctx context.Context, req *pb.CreateSlotRequest) (*pb.SlotResponse, error) {
	businessID, err := uuid.Parse(req.BusinessId)
	if err != nil {
		return nil, err
	}

	createReq := domain.CreateSlotRequest{
		BusinessID: businessID,
		Date:       req.Date,
		Slots:      req.Slots,
	}

	if err := s.useCase.CreateSlots(ctx, createReq); err != nil {
		return nil, err
	}

	return &pb.SlotResponse{
		BusinessId: req.BusinessId,
		Date:      req.Date,
	}, nil
}

func (s *Server) ListSlots(ctx context.Context, req *pb.SlotQuery) (*pb.SlotList, error) {
	businessID, err := uuid.Parse(req.BusinessId)
	if err != nil {
		return nil, err
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, err
	}

	slots, err := s.useCase.GetAvailableSlots(ctx, businessID, date)
	if err != nil {
		return nil, err
	}

	var response pb.SlotList
	for _, slot := range slots {
		response.Slots = append(response.Slots, &pb.SlotResponse{
			Id:         slot.ID.String(),
			BusinessId: slot.BusinessID.String(),
			Date:       slot.Date.Format("2006-01-02"),
			Time:       slot.Time.Format("15:04"),
			IsBooked:   slot.IsBooked,
		})
	}

	return &response, nil
}

func (s *Server) BookAppointment(ctx context.Context, req *pb.AppointmentRequest) (*pb.AppointmentResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	slotID, err := uuid.Parse(req.SlotId)
	if err != nil {
		return nil, err
	}

	bookReq := domain.BookAppointmentRequest{
		UserID: userID,
		SlotID: slotID,
	}

	appointment, err := s.useCase.BookAppointment(ctx, bookReq)
	if err != nil {
		return nil, err
	}

	return &pb.AppointmentResponse{
		Id:        appointment.ID.String(),
		UserId:    appointment.UserID.String(),
		SlotId:    appointment.SlotID.String(),
		CreatedAt: timestamppb.New(appointment.CreatedAt),
	}, nil
} 