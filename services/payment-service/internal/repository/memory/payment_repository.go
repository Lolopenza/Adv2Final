package memory

import (
	"sync"

	"payment-service/internal/domain"
)

type paymentRepository struct {
	payments map[string]*domain.Payment
	mutex    sync.RWMutex
}

// NewPaymentRepository creates a new in-memory payment repository for testing
func NewPaymentRepository() domain.PaymentRepository {
	return &paymentRepository{
		payments: make(map[string]*domain.Payment),
	}
}

func (r *paymentRepository) Create(payment *domain.Payment) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.payments[payment.ID] = payment
	return nil
}

func (r *paymentRepository) GetByID(id string) (*domain.Payment, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	payment, exists := r.payments[id]
	if !exists {
		return nil, domain.ErrPaymentNotFound
	}
	return payment, nil
}

func (r *paymentRepository) Update(payment *domain.Payment) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.payments[payment.ID] = payment
	return nil
}

func (r *paymentRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.payments, id)
	return nil
}

func (r *paymentRepository) List(customerEmail string, page, limit int32) ([]*domain.Payment, int32, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var result []*domain.Payment
	var total int32 = 0

	// Filter by customer email if provided
	for _, payment := range r.payments {
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
