package unit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"payment-service/internal/domain"
)

func TestPublishPaymentConfirmed(t *testing.T) {
	mockPublisher := new(MockPaymentEventPublisher)
	payment := &domain.Payment{
		ID:            "test-payment-id",
		Amount:        100.00,
		Currency:      "USD",
		Status:        "completed",
		CustomerEmail: "test@example.com",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	mockPublisher.On("PublishPaymentConfirmed", payment).Return(nil)
	// Execute
	err := mockPublisher.PublishPaymentConfirmed(payment)
	// Assert
	assert.NoError(t, err)
	mockPublisher.AssertExpectations(t)
}

func TestPublishPaymentFailed(t *testing.T) {
	mockPublisher := new(MockPaymentEventPublisher)
	payment := &domain.Payment{
		ID:            "test-payment-id",
		Amount:        100.00,
		Currency:      "USD",
		Status:        "failed",
		CustomerEmail: "test@example.com",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	mockPublisher.On("PublishPaymentFailed", payment).Return(nil)
	// Execute
	err := mockPublisher.PublishPaymentFailed(payment)
	// Assert
	assert.NoError(t, err)
	mockPublisher.AssertExpectations(t)
}

func TestPublishPaymentRefunded(t *testing.T) {
	mockPublisher := new(MockPaymentEventPublisher)
	payment := &domain.Payment{
		ID:            "test-payment-id",
		Amount:        100.00,
		Currency:      "USD",
		Status:        "refunded",
		CustomerEmail: "test@example.com",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	mockPublisher.On("PublishPaymentRefunded", payment).Return(nil)
	// Execute
	err := mockPublisher.PublishPaymentRefunded(payment)
	// Assert
	assert.NoError(t, err)
	mockPublisher.AssertExpectations(t)
}
