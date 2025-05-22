package grpc

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"payment-service/internal/domain"
	"payment-service/internal/usecase"
	pb "payment-service/proto"
)

type paymentServer struct {
	pb.UnimplementedPaymentServiceServer
	paymentUseCase domain.PaymentUseCase
}

func NewPaymentServer(paymentUseCase domain.PaymentUseCase) pb.PaymentServiceServer {
	return &paymentServer{
		paymentUseCase: paymentUseCase,
	}
}

func (s *paymentServer) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	payment := &domain.Payment{
		Amount:        req.Amount,
		Currency:      req.Currency,
		CustomerEmail: req.CustomerEmail,
		Description:   req.Description,
	}

	if err := s.paymentUseCase.InitiatePayment(payment); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreatePaymentResponse{
		PaymentId: payment.ID,
		Status:    string(payment.Status),
		CreatedAt: timestamppb.New(payment.CreatedAt),
	}, nil
}

func (s *paymentServer) ConfirmPayment(ctx context.Context, req *pb.ConfirmPaymentRequest) (*pb.ConfirmPaymentResponse, error) {
	if err := s.paymentUseCase.ConfirmPayment(req.PaymentId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	payment, err := s.paymentUseCase.GetPaymentStatus(req.PaymentId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ConfirmPaymentResponse{
		PaymentId: payment.ID,
		Status:    string(payment.Status),
		UpdatedAt: timestamppb.New(payment.UpdatedAt),
	}, nil
}

func (s *paymentServer) GetPayment(ctx context.Context, req *pb.GetPaymentRequest) (*pb.GetPaymentResponse, error) {
	payment, err := s.paymentUseCase.GetPaymentStatus(req.PaymentId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetPaymentResponse{
		PaymentId:     payment.ID,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		Status:        string(payment.Status),
		CustomerEmail: payment.CustomerEmail,
		Description:   payment.Description,
		CreatedAt:     timestamppb.New(payment.CreatedAt),
		UpdatedAt:     timestamppb.New(payment.UpdatedAt),
	}, nil
}

func (s *paymentServer) RefundPayment(ctx context.Context, req *pb.RefundPaymentRequest) (*pb.RefundPaymentResponse, error) {
	if err := s.paymentUseCase.RefundPayment(req.PaymentId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	payment, err := s.paymentUseCase.GetPaymentStatus(req.PaymentId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RefundPaymentResponse{
		PaymentId:  payment.ID,
		Status:     string(payment.Status),
		RefundedAt: timestamppb.New(payment.UpdatedAt),
	}, nil
}

func (s *paymentServer) ListPayments(ctx context.Context, req *pb.ListPaymentsRequest) (*pb.ListPaymentsResponse, error) {
	log.Printf("ListPayments called with request: %+v", req)

	// Set default values if not provided
	page := int32(1)
	limit := int32(10)

	// Only use values from request if they're valid
	if req != nil {
		if req.Page > 0 {
			page = req.Page
		}
		if req.Limit > 0 {
			limit = req.Limit
		}
	} else {
		log.Println("WARNING: Received nil request for ListPayments")
		req = &pb.ListPaymentsRequest{}
	}

	// Empty string is fine for customerEmail (no filter)
	customerEmail := ""
	if req != nil && req.CustomerEmail != "" {
		customerEmail = req.CustomerEmail
	}

	log.Printf("Querying with values: customerEmail=%q, page=%d, limit=%d", customerEmail, page, limit)

	payments, total, err := s.paymentUseCase.ListPayments(customerEmail, page, limit)
	if err != nil {
		log.Printf("Error in ListPayments: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Printf("Found %d payments out of %d total", len(payments), total)

	// Convert domain payments to proto payments
	protoPayments := make([]*pb.GetPaymentResponse, 0, len(payments))
	for _, payment := range payments {
		protoPayment := &pb.GetPaymentResponse{
			PaymentId:     payment.ID,
			Amount:        payment.Amount,
			Currency:      payment.Currency,
			Status:        string(payment.Status),
			CustomerEmail: payment.CustomerEmail,
			Description:   payment.Description,
			CreatedAt:     timestamppb.New(payment.CreatedAt),
			UpdatedAt:     timestamppb.New(payment.UpdatedAt),
		}
		protoPayments = append(protoPayments, protoPayment)
	}

	response := &pb.ListPaymentsResponse{
		Payments: protoPayments,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}

	log.Printf("Returning ListPaymentsResponse with %d payments", len(protoPayments))
	return response, nil
}

func (s *paymentServer) CancelPayment(ctx context.Context, req *pb.CancelPaymentRequest) (*pb.CancelPaymentResponse, error) {
	if err := s.paymentUseCase.CancelPayment(req.PaymentId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	payment, err := s.paymentUseCase.GetPaymentStatus(req.PaymentId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CancelPaymentResponse{
		PaymentId:   payment.ID,
		Status:      string(payment.Status),
		CancelledAt: timestamppb.New(payment.UpdatedAt),
	}, nil
}

func (s *paymentServer) GenerateInvoice(ctx context.Context, req *pb.GenerateInvoiceRequest) (*pb.GenerateInvoiceResponse, error) {
	log.Printf("GenerateInvoice called with request: %+v", req)

	invoiceURL, invoicePDF, err := s.paymentUseCase.GenerateInvoice(req.PaymentId)
	if err != nil {
		log.Printf("Error generating invoice: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Printf("Invoice generated successfully for payment %s, PDF size: %d bytes", req.PaymentId, len(invoicePDF))

	// Always send email, regardless of SendEmail flag
	log.Printf("Always sending invoice email for payment %s", req.PaymentId)
	payment, err := s.paymentUseCase.GetPaymentStatus(req.PaymentId)
	if err != nil {
		log.Printf("Error retrieving payment status: %v", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("invoice generated but failed to send email: %v", err))
	}

	log.Printf("Retrieved payment info: ID=%s, CustomerEmail=%s", payment.ID, payment.CustomerEmail)
	log.Printf("Sending invoice email for payment %s to %s", payment.ID, payment.CustomerEmail)

	if err := s.paymentUseCase.SendInvoiceEmail(payment, invoiceURL, invoicePDF); err != nil {
		log.Printf("Failed to send invoice email: %v", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("invoice generated but failed to send email: %v", err))
	}

	log.Printf("Successfully sent invoice email")

	return &pb.GenerateInvoiceResponse{
		PaymentId:  req.PaymentId,
		InvoiceUrl: invoiceURL,
		InvoicePdf: invoicePDF,
	}, nil
}

func (s *paymentServer) SendPaymentReminder(ctx context.Context, req *pb.SendPaymentReminderRequest) (*pb.SendPaymentReminderResponse, error) {
	if err := s.paymentUseCase.SendPaymentReminder(req.PaymentId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.SendPaymentReminderResponse{
		PaymentId: req.PaymentId,
		Success:   true,
		Message:   "Payment reminder sent successfully",
	}, nil
}

func (s *paymentServer) DeletePayment(ctx context.Context, req *pb.DeletePaymentRequest) (*pb.DeletePaymentResponse, error) {
	if err := s.paymentUseCase.DeletePayment(req.PaymentId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeletePaymentResponse{
		Success: true,
		Message: "Payment deleted successfully",
	}, nil
}

func (s *paymentServer) UpdatePayment(ctx context.Context, req *pb.UpdatePaymentRequest) (*pb.UpdatePaymentResponse, error) {
	log.Printf("UpdatePayment gRPC handler called with request: %+v", req)

	if req.PaymentId == "" {
		log.Printf("Missing payment ID in request")
		return nil, status.Error(codes.InvalidArgument, "payment ID is required")
	}

	log.Printf("Calling usecase UpdatePayment with id=%s, amount=%.2f, currency=%s, description=%s",
		req.PaymentId, req.Amount, req.Currency, req.Description)

	if err := s.paymentUseCase.UpdatePayment(req.PaymentId, req.Amount, req.Currency, req.Description); err != nil {
		log.Printf("Error in usecase UpdatePayment: %v", err)

		// Provide more specific error messages based on error type
		if err == usecase.ErrPaymentNotFound {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("payment with ID %s not found", req.PaymentId))
		} else if err == usecase.ErrInvalidStatus {
			return nil, status.Error(codes.FailedPrecondition, "only payments with PENDING status can be updated")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Printf("Payment updated successfully, fetching updated payment")

	payment, err := s.paymentUseCase.GetPaymentStatus(req.PaymentId)
	if err != nil {
		log.Printf("Error getting updated payment: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Printf("Returning updated payment: id=%s, status=%s, updatedAt=%v",
		payment.ID, payment.Status, payment.UpdatedAt)

	return &pb.UpdatePaymentResponse{
		PaymentId: payment.ID,
		Status:    string(payment.Status),
		UpdatedAt: timestamppb.New(payment.UpdatedAt),
	}, nil
}
