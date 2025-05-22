package nats

import (
	"encoding/json"

	"payment-service/internal/domain"

	"github.com/nats-io/nats.go"
)

const (
	PaymentConfirmedSubject = "payment.confirmed"
	PaymentFailedSubject    = "payment.failed"
	PaymentRefundedSubject  = "payment.refunded"
)

type eventPublisher struct {
	conn *nats.Conn
}

func NewEventPublisher(conn *nats.Conn) domain.PaymentEventPublisher {
	return &eventPublisher{conn: conn}
}

func (p *eventPublisher) PublishPaymentConfirmed(payment *domain.Payment) error {
	return p.publishEvent(PaymentConfirmedSubject, payment)
}

func (p *eventPublisher) PublishPaymentFailed(payment *domain.Payment) error {
	return p.publishEvent(PaymentFailedSubject, payment)
}

func (p *eventPublisher) PublishPaymentRefunded(payment *domain.Payment) error {
	return p.publishEvent(PaymentRefundedSubject, payment)
}

func (p *eventPublisher) publishEvent(subject string, payment *domain.Payment) error {
	data, err := json.Marshal(payment)
	if err != nil {
		return err
	}
	return p.conn.Publish(subject, data)
}
