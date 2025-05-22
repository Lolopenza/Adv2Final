package usecase

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"time"

	"payment-service/internal/domain"

	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
)

var (
	ErrPaymentNotFound = errors.New("payment not found")
	ErrInvalidStatus   = errors.New("invalid payment status")
)

type paymentUseCase struct {
	repo           domain.PaymentRepository
	cache          domain.PaymentCache
	eventPublisher domain.PaymentEventPublisher
	emailService   domain.EmailService
}

func NewPaymentUseCase(
	repo domain.PaymentRepository,
	cache domain.PaymentCache,
	eventPublisher domain.PaymentEventPublisher,
	emailService domain.EmailService,
) domain.PaymentUseCase {
	return &paymentUseCase{
		repo:           repo,
		cache:          cache,
		eventPublisher: eventPublisher,
		emailService:   emailService,
	}
}

func (uc *paymentUseCase) InitiatePayment(payment *domain.Payment) error {
	payment.ID = uuid.New().String()
	payment.Status = domain.PaymentStatusPending
	payment.CreatedAt = time.Now()
	payment.UpdatedAt = time.Now()

	if err := uc.repo.Create(payment); err != nil {
		return err
	}

	// Cache the payment for quick access
	cacheKey := GeneratePaymentCacheKey(payment.ID)
	if err := uc.cache.Set(cacheKey, payment, 24*time.Hour); err != nil {
		// Log cache error but don't fail the operation
		// TODO: Add proper logging
	}

	return nil
}

func (uc *paymentUseCase) ConfirmPayment(id string) error {
	payment, err := uc.getPayment(id)
	if err != nil {
		return err
	}

	if payment.Status != domain.PaymentStatusPending {
		return ErrInvalidStatus
	}

	payment.Status = domain.PaymentStatusCompleted
	payment.UpdatedAt = time.Now()

	if err := uc.repo.Update(payment); err != nil {
		return err
	}

	// Update cache
	cacheKey := GeneratePaymentCacheKey(payment.ID)
	if err := uc.cache.Set(cacheKey, payment, 24*time.Hour); err != nil {
		log.Printf("Warning: Failed to update payment in cache: %v", err)
	}

	// Publish event
	if err := uc.eventPublisher.PublishPaymentConfirmed(payment); err != nil {
		log.Printf("Warning: Failed to publish payment confirmation event: %v", err)
	}

	// Send confirmation email
	if err := uc.emailService.SendPaymentConfirmation(payment); err != nil {
		log.Printf("Warning: Failed to send payment confirmation email: %v", err)
	}

	// Automatically generate invoice for completed payment
	log.Printf("Payment %s confirmed, generating invoice", payment.ID)
	invoiceURL, invoicePDF, err := uc.GenerateInvoice(payment.ID)
	if err != nil {
		log.Printf("Warning: Failed to generate invoice for payment %s: %v", payment.ID, err)
	} else {
		// Send invoice email
		log.Printf("Successfully generated invoice, sending email to %s", payment.CustomerEmail)
		if err := uc.SendInvoiceEmail(payment, invoiceURL, invoicePDF); err != nil {
			log.Printf("Warning: Failed to send invoice email for payment %s: %v", payment.ID, err)
		} else {
			log.Printf("Successfully sent invoice email for completed payment %s", payment.ID)
		}
	}

	return nil
}

func (uc *paymentUseCase) GetPaymentStatus(id string) (*domain.Payment, error) {
	// Try to get from cache first
	cacheKey := GeneratePaymentCacheKey(id)
	if cached, err := uc.cache.Get(cacheKey); err == nil && cached != nil {
		return cached.(*domain.Payment), nil
	}

	// If not in cache, get from database
	payment, err := uc.getPayment(id)
	if err != nil {
		return nil, err
	}

	// Update cache
	if err := uc.cache.Set(cacheKey, payment, 24*time.Hour); err != nil {
		// Log cache error but don't fail the operation
		// TODO: Add proper logging
	}

	return payment, nil
}

func (uc *paymentUseCase) RefundPayment(id string) error {
	payment, err := uc.getPayment(id)
	if err != nil {
		return err
	}

	if payment.Status != domain.PaymentStatusCompleted {
		return ErrInvalidStatus
	}

	payment.Status = domain.PaymentStatusRefunded
	payment.UpdatedAt = time.Now()

	if err := uc.repo.Update(payment); err != nil {
		return err
	}

	// Update cache
	cacheKey := GeneratePaymentCacheKey(payment.ID)
	if err := uc.cache.Set(cacheKey, payment, 24*time.Hour); err != nil {
		// Log cache error but don't fail the operation
		// TODO: Add proper logging
	}

	// Publish event
	if err := uc.eventPublisher.PublishPaymentRefunded(payment); err != nil {
		// Log event publishing error but don't fail the operation
		// TODO: Add proper logging
	}

	// Send refund confirmation email
	if err := uc.emailService.SendRefundConfirmation(payment); err != nil {
		// Log email error but don't fail the operation
		// TODO: Add proper logging
	}

	return nil
}

func (uc *paymentUseCase) ListPayments(customerEmail string, page, limit int32) ([]*domain.Payment, int32, error) {
	log.Printf("UseCase ListPayments called with: customerEmail=%q, page=%d, limit=%d", customerEmail, page, limit)

	if page < 1 {
		page = 1
		log.Printf("Adjusted page to minimum value: %d", page)
	}
	if limit < 1 || limit > 100 {
		limit = 10 // Default limit
		log.Printf("Adjusted limit to default value: %d", limit)
	}

	log.Printf("Calling repository List method")
	payments, total, err := uc.repo.List(customerEmail, page, limit)
	if err != nil {
		log.Printf("Error in repository List method: %v", err)
		return nil, 0, err
	}

	log.Printf("Successfully retrieved %d payments (total: %d)", len(payments), total)
	return payments, total, nil
}

func (uc *paymentUseCase) CancelPayment(id string) error {
	payment, err := uc.getPayment(id)
	if err != nil {
		return err
	}

	if payment.Status != domain.PaymentStatusPending {
		return ErrInvalidStatus
	}

	payment.Status = domain.PaymentStatusFailed
	payment.UpdatedAt = time.Now()

	if err := uc.repo.Update(payment); err != nil {
		return err
	}

	// Update cache
	cacheKey := GeneratePaymentCacheKey(payment.ID)
	if err := uc.cache.Set(cacheKey, payment, 24*time.Hour); err != nil {
		// Log cache error but don't fail the operation
		// TODO: Add proper logging
	}

	// Publish event
	if err := uc.eventPublisher.PublishPaymentFailed(payment); err != nil {
		// Log event publishing error but don't fail the operation
		// TODO: Add proper logging
	}

	return nil
}

func (uc *paymentUseCase) GenerateInvoice(id string) (string, []byte, error) {
	payment, err := uc.getPayment(id)
	if err != nil {
		return "", nil, err
	}

	if payment.Status != domain.PaymentStatusCompleted {
		return "", nil, ErrInvalidStatus
	}

	// Generate a download URL (in production this might be a pre-signed URL or a route to a file server)
	invoiceURL := "https://api.example.com/invoices/" + payment.ID + ".pdf"

	// Create invoice PDF using gofpdf
	invoicePDF, err := generateProfessionalInvoice(payment)
	if err != nil {
		log.Printf("Error generating PDF invoice: %v", err)
		return "", nil, err
	}

	return invoiceURL, invoicePDF, nil
}

// generateProfessionalInvoice creates a well-formatted PDF invoice
func generateProfessionalInvoice(payment *domain.Payment) ([]byte, error) {
	// Create a new PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set fonts
	pdf.SetFont("Helvetica", "B", 16)
	pdf.SetTextColor(0, 51, 102) // Dark blue

	// Company Information
	pdf.Cell(40, 10, "SecurePayments Inc.")
	pdf.Ln(8)
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(80, 80, 80) // Dark gray
	pdf.Cell(40, 6, "200 Payment Plaza, Suite 300")
	pdf.Ln(6)
	pdf.Cell(40, 6, "Financial District, NY 10004")
	pdf.Ln(6)
	pdf.Cell(40, 6, "support@securepayments.com")
	pdf.Ln(6)
	pdf.Cell(40, 6, "+1 (800) PAY-SAFE")
	pdf.Ln(15)

	// Invoice Title and Number
	pdf.SetFont("Helvetica", "B", 14)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(0, 10, "INVOICE")
	pdf.Ln(6)
	pdf.SetFont("Helvetica", "", 10)
	pdf.Cell(0, 10, "Invoice Number: INV-"+payment.ID[:8])
	pdf.Ln(6)
	pdf.Cell(0, 10, "Date: "+payment.CreatedAt.Format("January 2, 2006"))
	pdf.Ln(15)

	// Customer Information
	pdf.SetFont("Helvetica", "B", 12)
	pdf.Cell(0, 10, "Billed To:")
	pdf.Ln(8)
	pdf.SetFont("Helvetica", "", 10)
	pdf.Cell(0, 6, payment.CustomerEmail)
	pdf.Ln(15)

	// Invoice Details Table
	pdf.SetFont("Helvetica", "B", 10)
	pdf.SetFillColor(240, 240, 240)

	// Table header
	pdf.CellFormat(100, 8, "Description", "1", 0, "", true, 0, "")
	pdf.CellFormat(30, 8, "Quantity", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Unit Price", "1", 0, "R", true, 0, "")
	pdf.CellFormat(30, 8, "Amount", "1", 1, "R", true, 0, "")

	// Table content
	pdf.SetFont("Helvetica", "", 10)
	pdf.CellFormat(100, 8, payment.Description, "1", 0, "", false, 0, "")
	pdf.CellFormat(30, 8, "1", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 8, fmt.Sprintf("%.2f %s", payment.Amount, payment.Currency), "1", 0, "R", false, 0, "")
	pdf.CellFormat(30, 8, fmt.Sprintf("%.2f %s", payment.Amount, payment.Currency), "1", 1, "R", false, 0, "")

	// Totals
	pdf.CellFormat(160, 8, "Subtotal:", "1", 0, "R", false, 0, "")
	pdf.CellFormat(30, 8, fmt.Sprintf("%.2f %s", payment.Amount, payment.Currency), "1", 1, "R", false, 0, "")

	pdf.CellFormat(160, 8, "Tax (0%):", "1", 0, "R", false, 0, "")
	pdf.CellFormat(30, 8, "0.00 "+payment.Currency, "1", 1, "R", false, 0, "")

	pdf.SetFont("Helvetica", "B", 10)
	pdf.CellFormat(160, 8, "Total:", "1", 0, "R", false, 0, "")
	pdf.CellFormat(30, 8, fmt.Sprintf("%.2f %s", payment.Amount, payment.Currency), "1", 1, "R", false, 0, "")

	// Payment Information
	pdf.Ln(10)
	pdf.SetFont("Helvetica", "B", 10)
	pdf.Cell(0, 8, "Payment Information")
	pdf.Ln(8)
	pdf.SetFont("Helvetica", "", 10)
	pdf.Cell(0, 6, "Payment ID: "+payment.ID)
	pdf.Ln(6)
	pdf.Cell(0, 6, "Status: "+string(payment.Status))
	pdf.Ln(6)
	pdf.Cell(0, 6, "Date: "+payment.UpdatedAt.Format("January 2, 2006 15:04:05"))

	// Thank you note
	pdf.Ln(15)
	pdf.Cell(0, 10, "Thank you for your business!")

	// Footer with Page number
	pdf.SetY(-15)
	pdf.SetFont("Helvetica", "I", 8)
	pdf.Cell(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()))

	// Convert to byte array
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (uc *paymentUseCase) SendPaymentReminder(id string) error {
	payment, err := uc.getPayment(id)
	if err != nil {
		return err
	}

	if payment.Status != domain.PaymentStatusPending {
		return ErrInvalidStatus
	}

	// Send the reminder email
	if err := uc.emailService.SendPaymentReminder(payment); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
		return err
	}

	return nil
}

func (uc *paymentUseCase) DeletePayment(id string) error {
	// Check if payment exists
	payment, err := uc.getPayment(id)
	if err != nil {
		return err
	}

	// Delete from repository
	if err := uc.repo.Delete(payment.ID); err != nil {
		return err
	}

	// Delete from cache
	cacheKey := GeneratePaymentCacheKey(payment.ID)
	if err := uc.cache.Delete(cacheKey); err != nil {
		// Log cache error but don't fail the operation
		// TODO: Add proper logging
	}

	return nil
}

func (uc *paymentUseCase) UpdatePayment(id string, amount float64, currency, description string) error {
	log.Printf("UpdatePayment called with id=%s, amount=%.2f, currency=%s, description=%s", id, amount, currency, description)

	payment, err := uc.getPayment(id)
	if err != nil {
		log.Printf("Error getting payment: %v", err)
		return err
	}

	// Can only update pending payments
	if payment.Status != domain.PaymentStatusPending {
		log.Printf("Cannot update payment with status %s, only PENDING payments can be updated", payment.Status)
		return ErrInvalidStatus
	}

	log.Printf("Original payment: amount=%.2f, currency=%s, description=%s", payment.Amount, payment.Currency, payment.Description)

	// Update payment fields
	payment.Amount = amount
	payment.Currency = currency
	payment.Description = description
	payment.UpdatedAt = time.Now()

	log.Printf("Updated payment: amount=%.2f, currency=%s, description=%s, updatedAt=%v", payment.Amount, payment.Currency, payment.Description, payment.UpdatedAt)

	// Save to repository
	if err := uc.repo.Update(payment); err != nil {
		log.Printf("Error updating payment in repository: %v", err)
		return err
	}
	log.Printf("Payment successfully updated in repository")

	// Update cache
	cacheKey := GeneratePaymentCacheKey(payment.ID)
	if err := uc.cache.Set(cacheKey, payment, 24*time.Hour); err != nil {
		// Log cache error but don't fail the operation
		log.Printf("Warning: Failed to update payment in cache: %v", err)
	} else {
		log.Printf("Payment successfully updated in cache")
	}

	return nil
}

func (uc *paymentUseCase) getPayment(id string) (*domain.Payment, error) {
	payment, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}
	return payment, nil
}

func GeneratePaymentCacheKey(paymentID string) string {
	return "payment:" + paymentID
}

func (uc *paymentUseCase) SendInvoiceEmail(payment *domain.Payment, invoiceURL string, invoicePDF []byte) error {
	// Add detailed logging
	log.Printf("SendInvoiceEmail called for payment ID: %s to email: %s", payment.ID, payment.CustomerEmail)
	log.Printf("Invoice URL: %s, PDF Size: %d bytes", invoiceURL, len(invoicePDF))

	// Send the invoice email
	if err := uc.emailService.SendInvoiceEmail(payment, invoiceURL, invoicePDF); err != nil {
		log.Printf("Error in emailService.SendInvoiceEmail: %v", err)
		return err
	}

	log.Printf("Successfully called emailService.SendInvoiceEmail for payment ID: %s", payment.ID)
	return nil
}
