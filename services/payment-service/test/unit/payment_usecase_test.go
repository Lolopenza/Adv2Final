package unit

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"payment-service/internal/domain"
	"payment-service/internal/usecase"
)

// Mock implementations
type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) Create(payment *domain.Payment) error {
	args := m.Called(payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) GetByID(id string) (*domain.Payment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payment), args.Error(1)
}

func (m *MockPaymentRepository) Update(payment *domain.Payment) error {
	args := m.Called(payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPaymentRepository) List(customerEmail string, page, limit int32) ([]*domain.Payment, int32, error) {
	args := m.Called(customerEmail, page, limit)
	return args.Get(0).([]*domain.Payment), args.Get(1).(int32), args.Error(2)
}

func (m *MockPaymentRepository) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

type MockPaymentCache struct {
	mock.Mock
}

func (m *MockPaymentCache) Set(key string, value interface{}, expiration time.Duration) error {
	args := m.Called(key, value, expiration)
	return args.Error(0)
}

func (m *MockPaymentCache) Get(key string) (interface{}, error) {
	args := m.Called(key)
	return args.Get(0), args.Error(1)
}

func (m *MockPaymentCache) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

type MockPaymentEventPublisher struct {
	mock.Mock
}

func (m *MockPaymentEventPublisher) PublishPaymentConfirmed(payment *domain.Payment) error {
	args := m.Called(payment)
	return args.Error(0)
}

func (m *MockPaymentEventPublisher) PublishPaymentFailed(payment *domain.Payment) error {
	args := m.Called(payment)
	return args.Error(0)
}

func (m *MockPaymentEventPublisher) PublishPaymentRefunded(payment *domain.Payment) error {
	args := m.Called(payment)
	return args.Error(0)
}

func (m *MockPaymentEventPublisher) PublishSubscriptionCreated(subscription *domain.Subscription) error {
	args := m.Called(subscription)
	return args.Error(0)
}

func (m *MockPaymentEventPublisher) PublishSubscriptionCancelled(subscription *domain.Subscription) error {
	args := m.Called(subscription)
	return args.Error(0)
}

func (m *MockPaymentEventPublisher) PublishSubscriptionRenewed(subscription *domain.Subscription) error {
	args := m.Called(subscription)
	return args.Error(0)
}

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendPaymentConfirmation(payment *domain.Payment) error {
	args := m.Called(payment)
	return args.Error(0)
}

func (m *MockEmailService) SendPaymentReceipt(payment *domain.Payment) error {
	args := m.Called(payment)
	return args.Error(0)
}

func (m *MockEmailService) SendRefundConfirmation(payment *domain.Payment) error {
	args := m.Called(payment)
	return args.Error(0)
}

func (m *MockEmailService) SendPaymentReminder(payment *domain.Payment) error {
	args := m.Called(payment)
	return args.Error(0)
}

func (m *MockEmailService) SendInvoiceEmail(payment *domain.Payment, invoiceURL string, invoicePDF []byte) error {
	args := m.Called(payment, invoiceURL, invoicePDF)
	return args.Error(0)
}

// Add subscription email methods
func (m *MockEmailService) SendSubscriptionConfirmation(customerEmail string) error {
	args := m.Called(customerEmail)
	return args.Error(0)
}

func (m *MockEmailService) SendCancellationEmail(customerEmail string) error {
	args := m.Called(customerEmail)
	return args.Error(0)
}

func (m *MockEmailService) SendRenewalEmail(customerEmail string) error {
	args := m.Called(customerEmail)
	return args.Error(0)
}

// Test cases
func TestInitiatePayment(t *testing.T) {
	// Setup
	mockRepo := new(MockPaymentRepository)
	mockCache := new(MockPaymentCache)
	mockPublisher := new(MockPaymentEventPublisher)
	mockEmail := new(MockEmailService)

	useCase := usecase.NewPaymentUseCase(mockRepo, mockCache, mockPublisher, mockEmail)

	// Test data
	payment := &domain.Payment{
		Amount:        100.00,
		Currency:      "USD",
		CustomerEmail: "test@example.com",
		Description:   "Test payment",
	}

	// Expectations
	mockRepo.On("Create", mock.AnythingOfType("*domain.Payment")).Return(nil)
	mockCache.On("Set", mock.AnythingOfType("string"), mock.AnythingOfType("*domain.Payment"), mock.AnythingOfType("time.Duration")).Return(nil)

	// Execute
	err := useCase.InitiatePayment(payment)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, payment.ID)
	assert.Equal(t, domain.PaymentStatusPending, payment.Status)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestConfirmPayment(t *testing.T) {
	// Setup
	mockRepo := new(MockPaymentRepository)
	mockCache := new(MockPaymentCache)
	mockPublisher := new(MockPaymentEventPublisher)
	mockEmail := new(MockEmailService)

	useCase := usecase.NewPaymentUseCase(mockRepo, mockCache, mockPublisher, mockEmail)

	// Test data
	paymentID := uuid.New().String()
	payment := &domain.Payment{
		ID:            paymentID,
		Amount:        100.00,
		Currency:      "USD",
		Status:        domain.PaymentStatusPending,
		CustomerEmail: "test@example.com",
		Description:   "Test payment",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Expectations
	// GetByID будет вызван дважды: первый раз в ConfirmPayment, второй раз в GenerateInvoice
	mockRepo.On("GetByID", paymentID).Return(payment, nil).Times(2)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Payment")).Return(nil)
	mockCache.On("Set", mock.AnythingOfType("string"), mock.AnythingOfType("*domain.Payment"), mock.AnythingOfType("time.Duration")).Return(nil)
	mockPublisher.On("PublishPaymentConfirmed", mock.AnythingOfType("*domain.Payment")).Return(nil)
	mockEmail.On("SendPaymentConfirmation", mock.AnythingOfType("*domain.Payment")).Return(nil)
	// Add expectation for SendInvoiceEmail
	mockEmail.On("SendInvoiceEmail", mock.AnythingOfType("*domain.Payment"), mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)

	// Execute
	err := useCase.ConfirmPayment(paymentID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
}

func TestRefundPayment(t *testing.T) {
	// Setup
	mockRepo := new(MockPaymentRepository)
	mockCache := new(MockPaymentCache)
	mockPublisher := new(MockPaymentEventPublisher)
	mockEmail := new(MockEmailService)

	useCase := usecase.NewPaymentUseCase(mockRepo, mockCache, mockPublisher, mockEmail)

	// Test data
	paymentID := uuid.New().String()
	payment := &domain.Payment{
		ID:            paymentID,
		Amount:        100.00,
		Currency:      "USD",
		Status:        domain.PaymentStatusCompleted,
		CustomerEmail: "test@example.com",
		Description:   "Test payment",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Expectations
	mockRepo.On("GetByID", paymentID).Return(payment, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Payment")).Return(nil)
	mockCache.On("Set", mock.AnythingOfType("string"), mock.AnythingOfType("*domain.Payment"), mock.AnythingOfType("time.Duration")).Return(nil)
	mockPublisher.On("PublishPaymentRefunded", mock.AnythingOfType("*domain.Payment")).Return(nil)
	mockEmail.On("SendRefundConfirmation", mock.AnythingOfType("*domain.Payment")).Return(nil)

	// Execute
	err := useCase.RefundPayment(paymentID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
}

func TestCancelPayment(t *testing.T) {
	// Setup
	mockRepo := new(MockPaymentRepository)
	mockCache := new(MockPaymentCache)
	mockPublisher := new(MockPaymentEventPublisher)
	mockEmail := new(MockEmailService)

	useCase := usecase.NewPaymentUseCase(mockRepo, mockCache, mockPublisher, mockEmail)

	// Test data
	paymentID := uuid.New().String()
	payment := &domain.Payment{
		ID:            paymentID,
		Amount:        100.00,
		Currency:      "USD",
		Status:        domain.PaymentStatusPending,
		CustomerEmail: "test@example.com",
		Description:   "Test payment",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Expectations
	mockRepo.On("GetByID", paymentID).Return(payment, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Payment")).Return(nil)
	mockCache.On("Set", mock.AnythingOfType("string"), mock.AnythingOfType("*domain.Payment"), mock.AnythingOfType("time.Duration")).Return(nil)
	mockPublisher.On("PublishPaymentFailed", mock.AnythingOfType("*domain.Payment")).Return(nil)

	// Execute
	err := useCase.CancelPayment(paymentID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestGenerateInvoice(t *testing.T) {
	// Setup
	mockRepo := new(MockPaymentRepository)
	mockCache := new(MockPaymentCache)
	mockPublisher := new(MockPaymentEventPublisher)
	mockEmail := new(MockEmailService)

	useCase := usecase.NewPaymentUseCase(mockRepo, mockCache, mockPublisher, mockEmail)

	// Test data
	paymentID := uuid.New().String()
	payment := &domain.Payment{
		ID:            paymentID,
		Amount:        100.00,
		Currency:      "USD",
		Status:        domain.PaymentStatusCompleted,
		CustomerEmail: "test@example.com",
		Description:   "Test payment",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Expectations
	mockRepo.On("GetByID", paymentID).Return(payment, nil)

	// Execute
	invoiceURL, invoicePDF, err := useCase.GenerateInvoice(paymentID)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, invoiceURL)
	assert.NotEmpty(t, invoicePDF)
	mockRepo.AssertExpectations(t)
}

func TestSendPaymentReminder(t *testing.T) {
	// Setup
	mockRepo := new(MockPaymentRepository)
	mockCache := new(MockPaymentCache)
	mockPublisher := new(MockPaymentEventPublisher)
	mockEmail := new(MockEmailService)

	useCase := usecase.NewPaymentUseCase(mockRepo, mockCache, mockPublisher, mockEmail)

	// Test data
	paymentID := uuid.New().String()
	payment := &domain.Payment{
		ID:            paymentID,
		Amount:        100.00,
		Currency:      "USD",
		Status:        domain.PaymentStatusPending,
		CustomerEmail: "test@example.com",
		Description:   "Test payment",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Expectations
	mockRepo.On("GetByID", paymentID).Return(payment, nil)
	mockEmail.On("SendPaymentReminder", mock.AnythingOfType("*domain.Payment")).Return(nil)

	// Execute
	err := useCase.SendPaymentReminder(paymentID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
}

func TestDeletePayment(t *testing.T) {
	// Setup
	mockRepo := new(MockPaymentRepository)
	mockCache := new(MockPaymentCache)
	mockPublisher := new(MockPaymentEventPublisher)
	mockEmail := new(MockEmailService)

	useCase := usecase.NewPaymentUseCase(mockRepo, mockCache, mockPublisher, mockEmail)

	// Test data
	paymentID := uuid.New().String()
	payment := &domain.Payment{
		ID:            paymentID,
		Amount:        100.00,
		Currency:      "USD",
		Status:        domain.PaymentStatusPending,
		CustomerEmail: "test@example.com",
		Description:   "Test payment",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Expectations
	mockRepo.On("GetByID", paymentID).Return(payment, nil)
	mockRepo.On("Delete", paymentID).Return(nil)
	mockCache.On("Delete", mock.AnythingOfType("string")).Return(nil)

	// Execute
	err := useCase.DeletePayment(paymentID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestUpdatePayment(t *testing.T) {
	// Setup
	mockRepo := new(MockPaymentRepository)
	mockCache := new(MockPaymentCache)
	mockPublisher := new(MockPaymentEventPublisher)
	mockEmail := new(MockEmailService)

	useCase := usecase.NewPaymentUseCase(mockRepo, mockCache, mockPublisher, mockEmail)

	// Test data
	paymentID := uuid.New().String()
	payment := &domain.Payment{
		ID:            paymentID,
		Amount:        100.00,
		Currency:      "USD",
		Status:        domain.PaymentStatusPending,
		CustomerEmail: "test@example.com",
		Description:   "Test payment",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	newAmount := 150.00
	newCurrency := "EUR"
	newDescription := "Updated test payment"

	// Expectations
	mockRepo.On("GetByID", paymentID).Return(payment, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Payment")).Return(nil)
	mockCache.On("Set", mock.AnythingOfType("string"), mock.AnythingOfType("*domain.Payment"), mock.AnythingOfType("time.Duration")).Return(nil)

	// Execute
	err := useCase.UpdatePayment(paymentID, newAmount, newCurrency, newDescription)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestSendInvoiceEmailUseCase(t *testing.T) {
	// Setup
	mockRepo := new(MockPaymentRepository)
	mockCache := new(MockPaymentCache)
	mockPublisher := new(MockPaymentEventPublisher)
	mockEmail := new(MockEmailService)

	useCase := usecase.NewPaymentUseCase(mockRepo, mockCache, mockPublisher, mockEmail)

	// Test data
	paymentID := uuid.New().String()
	payment := &domain.Payment{
		ID:            paymentID,
		Amount:        100.00,
		Currency:      "USD",
		Status:        domain.PaymentStatusCompleted,
		CustomerEmail: "test@example.com",
		Description:   "Test payment",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	invoiceURL := "https://example.com/invoices/" + paymentID + ".pdf"
	invoicePDF := []byte("test invoice content")

	// Expectations
	mockEmail.On("SendInvoiceEmail", payment, invoiceURL, invoicePDF).Return(nil)

	// Execute
	err := useCase.SendInvoiceEmail(payment, invoiceURL, invoicePDF)

	// Assert
	assert.NoError(t, err)
	mockEmail.AssertExpectations(t)
}
