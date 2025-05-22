package email

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	emailDelivery "payment-service/internal/delivery/email"
	"payment-service/internal/domain"
	"payment-service/internal/usecase"
	"payment-service/test/util"
)

func TestSendPaymentReminderEmail(t *testing.T) {
	// Start SMTP capture server
	smtpServer, err := util.NewSMTPCapture(2525)
	if err != nil {
		t.Fatalf("Failed to start SMTP capture server: %v", err)
	}

	smtpServer.Start()
	defer smtpServer.Stop()

	// Get host and port from SMTP server
	smtpHost, smtpPort := smtpServer.GetConfig()

	// Create test payment
	paymentID := uuid.New().String()
	testPayment := &domain.Payment{
		ID:            paymentID,
		Amount:        250.50,
		Currency:      "USD",
		Status:        domain.PaymentStatusPending,
		CustomerEmail: "customer@example.com",
		Description:   "Monthly subscription",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Create repo and add payment
	mockRepo := &testPaymentRepository{
		payments: map[string]*domain.Payment{
			paymentID: testPayment,
		},
	}

	// Create real email service with our SMTP capture server
	emailService := emailDelivery.NewEmailService(
		"test@service.com",
		"testpass",
		smtpHost,
		smtpPort,
	)

	// Create use case with real email service
	paymentUseCase := usecase.NewPaymentUseCase(
		mockRepo,
		&testPaymentCache{},
		&testPaymentEventPublisher{},
		emailService,
	)

	// Send reminder
	err = paymentUseCase.SendPaymentReminder(paymentID)
	require.NoError(t, err)

	// Wait a moment for email to be processed
	time.Sleep(500 * time.Millisecond)

	// Check captured emails
	messages := smtpServer.GetMessages()
	require.Equal(t, 1, len(messages), "Expected 1 email to be sent")

	msg := messages[0]

	// Verify email details
	assert.Equal(t, "test@service.com", msg.From, "From address doesn't match")
	assert.Equal(t, []string{"customer@example.com"}, msg.To, "To address doesn't match")
	assert.Equal(t, "Payment Reminder", msg.Subject, "Subject doesn't match")

	// Verify email content
	emailBody := string(msg.Body)

	// Check for essential payment details
	assert.Contains(t, emailBody, paymentID, "Email should contain payment ID")
	assert.Contains(t, emailBody, "250.50", "Email should contain payment amount")
	assert.Contains(t, emailBody, "USD", "Email should contain currency")
	assert.Contains(t, emailBody, "Monthly subscription", "Email should contain description")

	// Check for key phrases that indicate good format and clear instructions
	assert.Contains(t, emailBody, "Dear Customer", "Email should have a professional greeting")
	assert.Contains(t, emailBody, "reminder", "Email should indicate it's a reminder")
	assert.Contains(t, emailBody, "pending", "Email should indicate payment is pending")
	assert.Contains(t, emailBody, "earliest convenience", "Email should have clear instructions")
	assert.Contains(t, emailBody, "Thank you", "Email should have a professional closing")
}

// Test implementations

type testPaymentRepository struct {
	payments map[string]*domain.Payment
}

func (r *testPaymentRepository) Create(payment *domain.Payment) error {
	r.payments[payment.ID] = payment
	return nil
}

func (r *testPaymentRepository) GetByID(id string) (*domain.Payment, error) {
	payment, exists := r.payments[id]
	if !exists {
		return nil, domain.ErrPaymentNotFound
	}
	return payment, nil
}

func (r *testPaymentRepository) Update(payment *domain.Payment) error {
	r.payments[payment.ID] = payment
	return nil
}

func (r *testPaymentRepository) Delete(id string) error {
	delete(r.payments, id)
	return nil
}

func (r *testPaymentRepository) List(customerEmail string, page, limit int32) ([]*domain.Payment, int32, error) {
	var result []*domain.Payment
	var total int32 = 0

	for _, payment := range r.payments {
		if customerEmail == "" || payment.CustomerEmail == customerEmail {
			result = append(result, payment)
			total++
		}
	}

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

type testPaymentCache struct{}

func (c *testPaymentCache) Set(key string, value interface{}, expiration time.Duration) error {
	return nil
}

func (c *testPaymentCache) Get(key string) (interface{}, error) {
	return nil, nil
}

func (c *testPaymentCache) Delete(key string) error {
	return nil
}

type testPaymentEventPublisher struct{}

func (p *testPaymentEventPublisher) PublishPaymentConfirmed(payment *domain.Payment) error {
	return nil
}

func (p *testPaymentEventPublisher) PublishPaymentFailed(payment *domain.Payment) error {
	return nil
}

func (p *testPaymentEventPublisher) PublishPaymentRefunded(payment *domain.Payment) error {
	return nil
}
