package unit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"payment-service/internal/domain"
)

func TestSendPaymentConfirmation(t *testing.T) {
	mockEmail := new(MockEmailService)
	payment := &domain.Payment{
		ID:            "test-payment-id",
		Amount:        100.00,
		Currency:      "USD",
		Status:        "completed",
		CustomerEmail: "test@example.com",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	mockEmail.On("SendPaymentConfirmation", payment).Return(nil)
	// Execute
	err := mockEmail.SendPaymentConfirmation(payment)
	// Assert
	assert.NoError(t, err)
	mockEmail.AssertExpectations(t)
}

func TestSendPaymentFailed(t *testing.T) {
	mockEmail := new(MockEmailService)
	payment := &domain.Payment{
		ID:            "test-payment-id",
		Amount:        100.00,
		Currency:      "USD",
		Status:        "failed",
		CustomerEmail: "test@example.com",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	mockEmail.On("SendPaymentReceipt", payment).Return(nil)
	// Execute
	err := mockEmail.SendPaymentReceipt(payment)
	// Assert
	assert.NoError(t, err)
	mockEmail.AssertExpectations(t)
}

func TestSendPaymentRefunded(t *testing.T) {
	mockEmail := new(MockEmailService)
	payment := &domain.Payment{
		ID:            "test-payment-id",
		Amount:        100.00,
		Currency:      "USD",
		Status:        "refunded",
		CustomerEmail: "test@example.com",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	mockEmail.On("SendRefundConfirmation", payment).Return(nil)
	// Execute
	err := mockEmail.SendRefundConfirmation(payment)
	// Assert
	assert.NoError(t, err)
	mockEmail.AssertExpectations(t)
}

func TestSendInvoiceEmail(t *testing.T) {
	mockEmail := new(MockEmailService)
	payment := &domain.Payment{
		ID:            "test-payment-id",
		Amount:        100.00,
		Currency:      "USD",
		Status:        "completed",
		CustomerEmail: "test@example.com",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	invoiceURL := "https://example.com/invoices/test-payment-id.pdf"
	invoicePDF := []byte("test invoice content")

	mockEmail.On("SendInvoiceEmail", payment, invoiceURL, invoicePDF).Return(nil)

	// Execute
	err := mockEmail.SendInvoiceEmail(payment, invoiceURL, invoicePDF)

	// Assert
	assert.NoError(t, err)
	mockEmail.AssertExpectations(t)
}
