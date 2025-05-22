package domain

import (
	"errors"
	"time"
)

// Common errors
var (
	ErrPaymentNotFound = errors.New("payment not found")
	ErrInvalidPayment  = errors.New("invalid payment data")
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusCompleted PaymentStatus = "COMPLETED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
	PaymentStatusRefunded  PaymentStatus = "REFUNDED"
)

type Payment struct {
	ID            string
	Amount        float64
	Currency      string
	Status        PaymentStatus
	CustomerEmail string
	Description   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type PaymentRepository interface {
	Create(payment *Payment) error
	GetByID(id string) (*Payment, error)
	Update(payment *Payment) error
	Delete(id string) error
	List(customerEmail string, page, limit int32) ([]*Payment, int32, error)
}

type PaymentUseCase interface {
	InitiatePayment(payment *Payment) error
	ConfirmPayment(id string) error
	GetPaymentStatus(id string) (*Payment, error)
	RefundPayment(id string) error
	ListPayments(customerEmail string, page, limit int32) ([]*Payment, int32, error)
	CancelPayment(id string) error
	GenerateInvoice(id string) (string, []byte, error)
	SendPaymentReminder(id string) error
	DeletePayment(id string) error
	UpdatePayment(id string, amount float64, currency, description string) error
	SendInvoiceEmail(payment *Payment, invoiceURL string, invoicePDF []byte) error
}

type PaymentCache interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (interface{}, error)
	Delete(key string) error
}

type PaymentEventPublisher interface {
	PublishPaymentConfirmed(payment *Payment) error
	PublishPaymentFailed(payment *Payment) error
	PublishPaymentRefunded(payment *Payment) error
}

type EmailService interface {
	SendPaymentConfirmation(payment *Payment) error
	SendPaymentReceipt(payment *Payment) error
	SendRefundConfirmation(payment *Payment) error
	SendPaymentReminder(payment *Payment) error
	SendInvoiceEmail(payment *Payment, invoiceURL string, invoicePDF []byte) error

	// Subscription-related emails
	SendCancellationEmail(customerEmail string) error
	SendRenewalEmail(customerEmail string) error
}
