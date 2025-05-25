package email

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"path/filepath"

	"payment-service/internal/domain"
)

type emailService struct {
	from     string
	password string
	smtpHost string
	smtpPort string
}

func NewEmailService(from, password, smtpHost, smtpPort string) domain.EmailService {
	return &emailService{
		from:     from,
		password: password,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
	}
}

func (s *emailService) SendPaymentConfirmation(payment *domain.Payment) error {
	subject := "Payment Confirmation"
	body := fmt.Sprintf(`
		Dear Customer,
		
		Your payment of %f %s has been confirmed.
		Payment ID: %s
		Description: %s
		
		Thank you for your business!
	`, payment.Amount, payment.Currency, payment.ID, payment.Description)

	return s.sendEmail(payment.CustomerEmail, subject, body)
}

func (s *emailService) SendPaymentReceipt(payment *domain.Payment) error {
	subject := "Payment Receipt"
	body := fmt.Sprintf(`
		Dear Customer,
		
		Here is your payment receipt:
		
		Payment ID: %s
		Amount: %f %s
		Date: %s
		Description: %s
		
		Thank you for your business!
	`, payment.ID, payment.Amount, payment.Currency, payment.CreatedAt.Format("2006-01-02 15:04:05"), payment.Description)

	return s.sendEmail(payment.CustomerEmail, subject, body)
}

func (s *emailService) SendRefundConfirmation(payment *domain.Payment) error {
	subject := "Payment Refund Confirmation"
	body := fmt.Sprintf(`
		Dear Customer,
		
		Your refund for payment %s has been processed.
		Amount: %f %s
		Date: %s
		
		The refunded amount will be credited to your original payment method.
		
		Thank you for your understanding.
	`, payment.ID, payment.Amount, payment.Currency, payment.UpdatedAt.Format("2006-01-02 15:04:05"))

	return s.sendEmail(payment.CustomerEmail, subject, body)
}

func (s *emailService) SendPaymentReminder(payment *domain.Payment) error {
	subject := "Payment Reminder"
	body := fmt.Sprintf(`
		Dear Customer,
		
		This is a friendly reminder about your pending payment:
		
		Payment ID: %s
		Amount: %.2f %s
		Description: %s
		Created: %s
		
		Please complete this payment at your earliest convenience.
		
		If you've already made this payment, please disregard this message.
		
		Thank you for your business!
	`, payment.ID, payment.Amount, payment.Currency, payment.Description, payment.CreatedAt.Format("2006-01-02 15:04:05"))

	return s.sendEmail(payment.CustomerEmail, subject, body)
}

func (s *emailService) SendInvoiceEmail(payment *domain.Payment, invoiceURL string, invoicePDF []byte) error {
	log.Printf("Sending invoice email to %s for payment %s", payment.CustomerEmail, payment.ID)

	// Create buffer for MIME message
	var buffer bytes.Buffer

	// Create multipart writer
	writer := multipart.NewWriter(&buffer)

	// Set headers
	boundary := writer.Boundary()

	// Write email headers
	headers := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: Your Payment Invoice\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: multipart/mixed; boundary=%s\r\n\r\n",
		s.from, payment.CustomerEmail, boundary)
	buffer.WriteString(headers)

	// Add text part
	textPart, _ := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type":        {"text/plain; charset=UTF-8"},
		"Content-Disposition": {"inline"},
	})

	body := fmt.Sprintf(`
Dear Customer,

Thank you for your payment.

Please find attached your invoice for payment:
Payment ID: %s
Amount: %.2f %s
Description: %s
Date: %s

You can also access your invoice online at: %s

Thank you for your business!
`, payment.ID, payment.Amount, payment.Currency, payment.Description,
		payment.CreatedAt.Format("2006-01-02 15:04:05"), invoiceURL)

	textPart.Write([]byte(body))

	// Add attachment
	if len(invoicePDF) > 0 {
		attachmentFilename := fmt.Sprintf("invoice-%s.pdf", payment.ID)

		h := make(textproto.MIMEHeader)
		h.Set("Content-Type", "application/pdf")
		h.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(attachmentFilename)))
		h.Set("Content-Transfer-Encoding", "base64")

		attachmentPart, _ := writer.CreatePart(h)
		encoder := base64.NewEncoder(base64.StdEncoding, attachmentPart)
		encoder.Write(invoicePDF)
		encoder.Close()
	}

	// Close the writer
	writer.Close()

	// Send the email with attachment
	auth := smtp.PlainAuth("", s.from, s.password, s.smtpHost)
	addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)

	log.Printf("Connecting to SMTP server at %s", addr)

	err := smtp.SendMail(addr, auth, s.from, []string{payment.CustomerEmail}, buffer.Bytes())
	if err != nil {
		log.Printf("Error sending invoice email: %v", err)
		return err
	}

	log.Printf("Successfully sent invoice email to %s", payment.CustomerEmail)
	return nil
}

func (s *emailService) sendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.from, s.password, s.smtpHost)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", s.from, to, subject, body)

	addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)
	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg))
}

func (s *emailService) SendCancellationEmail(customerEmail string) error {
	subject := "Your subscription has been cancelled"
	body := "Your subscription has been cancelled. We're sorry to see you go!"

	return s.sendEmail(customerEmail, subject, body)
}

func (s *emailService) SendRenewalEmail(customerEmail string) error {
	subject := "Your subscription has been renewed"
	body := "Your subscription has been renewed. Thank you for your continued support!"

	return s.sendEmail(customerEmail, subject, body)
}
