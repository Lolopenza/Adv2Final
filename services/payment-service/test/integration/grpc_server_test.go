package integration

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	grpcDelivery "payment-service/internal/delivery/grpc"
	"payment-service/internal/domain"
	"payment-service/internal/usecase"
	pb "payment-service/proto"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

// Mock repository for integration test
var testPaymentRepo *mockPaymentRepository

type mockPaymentRepository struct {
	payments map[string]*domain.Payment
}

func newMockPaymentRepository() *mockPaymentRepository {
	return &mockPaymentRepository{payments: make(map[string]*domain.Payment)}
}

func (m *mockPaymentRepository) Create(payment *domain.Payment) error {
	m.payments[payment.ID] = payment
	return nil
}
func (m *mockPaymentRepository) GetByID(id string) (*domain.Payment, error) {
	p, ok := m.payments[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}
func (m *mockPaymentRepository) Update(payment *domain.Payment) error {
	m.payments[payment.ID] = payment
	return nil
}
func (m *mockPaymentRepository) Delete(id string) error {
	delete(m.payments, id)
	return nil
}

func (m *mockPaymentRepository) List(customerEmail string, page, limit int32) ([]*domain.Payment, int32, error) {
	var result []*domain.Payment
	var total int32 = 0

	// Filter by customer email if provided
	for _, payment := range m.payments {
		if customerEmail == "" || payment.CustomerEmail == customerEmail {
			result = append(result, payment)
			total++
		}
	}

	// Basic pagination (not very efficient for large datasets but works for tests)
	start := (page - 1) * limit
	end := start + limit
	if start >= int32(len(result)) {
		return []*domain.Payment{}, total, nil
	}
	if end > int32(len(result)) {
		end = int32(len(result))
	}

	return result[start:end], total, nil
}

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()

	// Initialize test database
	_, err := initTestDB()
	if err != nil {
		panic(err)
	}

	// Initialize repositories and services
	testPaymentRepo = newMockPaymentRepository()
	mockCache := &mockPaymentCache{}
	mockEventPublisher := &mockPaymentEventPublisher{}
	mockEmailService := &grpcMockEmailService{}
	paymentUseCase := usecase.NewPaymentUseCase(testPaymentRepo, mockCache, mockEventPublisher, mockEmailService)

	// Register gRPC server
	pb.RegisterPaymentServiceServer(s, grpcDelivery.NewPaymentServer(paymentUseCase))

	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

func initTestDB() (*gorm.DB, error) {
	dsn := "host=localhost port=5432 user=postgres password=Fydfhniga2 dbname=payment_service sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Run migrations
	// TODO: Add migration logic here

	return db, nil
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

// Mock cache for integration test
type mockPaymentCache struct{}

func (m *mockPaymentCache) Set(key string, value interface{}, expiration time.Duration) error {
	return nil
}
func (m *mockPaymentCache) Get(key string) (interface{}, error) { return nil, nil }
func (m *mockPaymentCache) Delete(key string) error             { return nil }

// Mock event publisher for integration test
type mockPaymentEventPublisher struct{}

func (m *mockPaymentEventPublisher) PublishPaymentCreated(payment *domain.Payment) error { return nil }
func (m *mockPaymentEventPublisher) PublishPaymentConfirmed(payment *domain.Payment) error {
	return nil
}
func (m *mockPaymentEventPublisher) PublishPaymentRefunded(payment *domain.Payment) error { return nil }
func (m *mockPaymentEventPublisher) PublishPaymentFailed(payment *domain.Payment) error   { return nil }

// Mock email service for integration test
type grpcMockEmailService struct{}

func (m *grpcMockEmailService) SendPaymentConfirmation(payment *domain.Payment) error { return nil }
func (m *grpcMockEmailService) SendRefundConfirmation(payment *domain.Payment) error  { return nil }
func (m *grpcMockEmailService) SendPaymentReceipt(payment *domain.Payment) error      { return nil }
func (m *grpcMockEmailService) SendPaymentReminder(payment *domain.Payment) error     { return nil }
func (m *grpcMockEmailService) SendInvoiceEmail(payment *domain.Payment, invoiceURL string, invoicePDF []byte) error {
	return nil
}
func (m *grpcMockEmailService) SendCancellationEmail(customerEmail string) error { return nil }
func (m *grpcMockEmailService) SendRenewalEmail(customerEmail string) error      { return nil }

func TestCreatePayment(t *testing.T) {
	ctx := context.Background()
	// Используем grpc.DialContext, который помечен как устаревший, но будет поддерживаться в 1.x версиях
	// TODO: обновить на более новые методы, когда они станут стабильными
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()
	client := pb.NewPaymentServiceClient(conn)

	req := &pb.CreatePaymentRequest{
		Amount:        100.00,
		Currency:      "USD",
		CustomerEmail: "test@example.com",
		Description:   "Test payment",
	}
	resp, err := client.CreatePayment(ctx, req)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.PaymentId)
	assert.Equal(t, "PENDING", resp.Status)
}

func TestConfirmPayment(t *testing.T) {
	ctx := context.Background()
	// Используем grpc.DialContext, который помечен как устаревший, но будет поддерживаться в 1.x версиях
	// TODO: обновить на более новые методы, когда они станут стабильными
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()
	client := pb.NewPaymentServiceClient(conn)

	// First create a payment
	createReq := &pb.CreatePaymentRequest{
		Amount:        100.00,
		Currency:      "USD",
		CustomerEmail: "test@example.com",
		Description:   "Test payment",
	}
	createResp, err := client.CreatePayment(ctx, createReq)
	require.NoError(t, err)

	// Now confirm it
	confirmReq := &pb.ConfirmPaymentRequest{
		PaymentId: createResp.PaymentId,
	}
	confirmResp, err := client.ConfirmPayment(ctx, confirmReq)
	require.NoError(t, err)
	assert.Equal(t, createResp.PaymentId, confirmResp.PaymentId)
	assert.Equal(t, "COMPLETED", confirmResp.Status)
}

func TestSendPaymentReminder(t *testing.T) {
	ctx := context.Background()
	// Используем grpc.DialContext, который помечен как устаревший, но будет поддерживаться в 1.x версиях
	// TODO: обновить на более новые методы, когда они станут стабильными
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()
	client := pb.NewPaymentServiceClient(conn)

	// First create a payment
	createReq := &pb.CreatePaymentRequest{
		Amount:        150.00,
		Currency:      "USD",
		CustomerEmail: "customer@example.com",
		Description:   "Subscription payment",
	}
	createResp, err := client.CreatePayment(ctx, createReq)
	require.NoError(t, err)

	// Send payment reminder
	reminderReq := &pb.SendPaymentReminderRequest{
		PaymentId: createResp.PaymentId,
	}
	reminderResp, err := client.SendPaymentReminder(ctx, reminderReq)
	require.NoError(t, err)
	assert.Equal(t, createResp.PaymentId, reminderResp.PaymentId)
	assert.True(t, reminderResp.Success)
	assert.Equal(t, "Payment reminder sent successfully", reminderResp.Message)

	// Test with non-existent payment ID
	badReminderReq := &pb.SendPaymentReminderRequest{
		PaymentId: "non-existent-id",
	}
	_, err = client.SendPaymentReminder(ctx, badReminderReq)
	assert.Error(t, err, "Should error with non-existent payment ID")

	// Create a payment and confirm it (so it's no longer pending)
	createReq2 := &pb.CreatePaymentRequest{
		Amount:        200.00,
		Currency:      "USD",
		CustomerEmail: "customer2@example.com",
		Description:   "Completed payment",
	}
	createResp2, err := client.CreatePayment(ctx, createReq2)
	require.NoError(t, err)

	// Confirm the payment
	confirmReq := &pb.ConfirmPaymentRequest{
		PaymentId: createResp2.PaymentId,
	}
	_, err = client.ConfirmPayment(ctx, confirmReq)
	require.NoError(t, err)

	// Try to send reminder for completed payment
	completedReminderReq := &pb.SendPaymentReminderRequest{
		PaymentId: createResp2.PaymentId,
	}
	_, err = client.SendPaymentReminder(ctx, completedReminderReq)
	assert.Error(t, err, "Should error with completed payment")
}
