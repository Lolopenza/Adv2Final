package usecase

import (
	"errors"
	"log"
	"time"

	"payment-service/internal/domain"
	"payment-service/internal/repository/cache"

	"github.com/google/uuid"
)

var (
	ErrSubscriptionNotFound = errors.New("subscription not found")
	// Use the existing ErrInvalidStatus from payment_usecase.go
)

type subscriptionUseCase struct {
	repo           domain.SubscriptionRepository
	cache          domain.SubscriptionCache
	eventPublisher domain.PaymentEventPublisher
	emailService   domain.EmailService
}

func NewSubscriptionUseCase(
	repo domain.SubscriptionRepository,
	cache domain.SubscriptionCache,
	eventPublisher domain.PaymentEventPublisher,
	emailService domain.EmailService,
) domain.SubscriptionUseCase {
	return &subscriptionUseCase{
		repo:           repo,
		cache:          cache,
		eventPublisher: eventPublisher,
		emailService:   emailService,
	}
}

func (uc *subscriptionUseCase) CreateSubscription(subscription *domain.Subscription) error {
	// Generate ID if not provided
	if subscription.ID == "" {
		subscription.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	subscription.CreatedAt = now
	subscription.UpdatedAt = now

	// Default status to ACTIVE if not set
	if subscription.Status == "" {
		subscription.Status = domain.SubscriptionStatusActive
	}

	// Create payment record
	payment := &domain.Payment{
		ID:            uuid.New().String(),
		Amount:        subscription.Price,
		Currency:      subscription.Currency,
		Status:        domain.PaymentStatusPending,
		CustomerEmail: subscription.CustomerEmail,
		Description:   "Subscription to " + subscription.PlanName,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Just use regular create method since CreateWithPayment would need a type cast
	if err := uc.repo.Create(subscription); err != nil {
		return err
	}

	// Here we would typically create the payment record via a transaction
	// For demonstration purposes, we're just logging it
	log.Printf("Created subscription %s with payment details: ID=%s, Amount=%v %s, Status=%s, Customer=%s, Description=%s, Created=%s, Updated=%s",
		subscription.ID,
		payment.ID,
		payment.Amount,
		payment.Currency,
		payment.Status,
		payment.CustomerEmail,
		payment.Description,
		payment.CreatedAt.Format(time.RFC3339),
		payment.UpdatedAt.Format(time.RFC3339))

	// Cache the subscription
	cacheKey := cache.GenerateSubscriptionCacheKey(subscription.ID)
	if err := uc.cache.Set(cacheKey, subscription, 24*time.Hour); err != nil {
		// Log cache error but don't fail the operation
		log.Printf("Warning: Failed to cache subscription: %v", err)
	}

	return nil
}

func (uc *subscriptionUseCase) GetSubscription(id string) (*domain.Subscription, error) {
	// Try to get from cache first
	cacheKey := cache.GenerateSubscriptionCacheKey(id)
	if cached, err := uc.cache.Get(cacheKey); err == nil && cached != nil {
		subscription, ok := cached.(*domain.Subscription)
		if ok {
			return subscription, nil
		}
		// If type assertion fails, log and continue to get from DB
		log.Printf("Warning: Cache returned invalid type for subscription: %T", cached)
	}

	// If not in cache, get from database
	subscription, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if subscription == nil {
		return nil, ErrSubscriptionNotFound
	}

	// Update cache
	if err := uc.cache.Set(cacheKey, subscription, 24*time.Hour); err != nil {
		// Log cache error but don't fail the operation
		log.Printf("Warning: Failed to cache subscription: %v", err)
	}

	return subscription, nil
}

func (uc *subscriptionUseCase) CancelSubscription(id string) error {
	subscription, err := uc.GetSubscription(id)
	if err != nil {
		return err
	}

	if subscription.Status != domain.SubscriptionStatusActive {
		return ErrInvalidStatus
	}

	subscription.Status = domain.SubscriptionStatusCancelled
	subscription.UpdatedAt = time.Now()

	if err := uc.repo.Update(subscription); err != nil {
		return err
	}

	// Update cache
	cacheKey := cache.GenerateSubscriptionCacheKey(id)
	if err := uc.cache.Set(cacheKey, subscription, 24*time.Hour); err != nil {
		log.Printf("Warning: Failed to update subscription in cache: %v", err)
	}

	// Send cancellation email
	if err := uc.emailService.SendCancellationEmail(subscription.CustomerEmail); err != nil {
		log.Printf("Warning: Failed to send cancellation email: %v", err)
	}

	return nil
}

func (uc *subscriptionUseCase) RenewSubscription(id string) error {
	subscription, err := uc.GetSubscription(id)
	if err != nil {
		return err
	}

	// Only allow renewing active subscriptions
	if subscription.Status != domain.SubscriptionStatusActive {
		return ErrInvalidStatus
	}

	// Update end date - typically 1 month or 1 year from now
	now := time.Now()
	subscription.StartDate = now
	subscription.EndDate = now.AddDate(0, 1, 0) // Add 1 month
	subscription.UpdatedAt = now

	// Log the renewal payment (we would create it in a real implementation)
	log.Printf("Creating renewal payment for subscription %s", subscription.ID)

	// Update subscription
	if err := uc.repo.Update(subscription); err != nil {
		return err
	}

	// Update cache
	cacheKey := cache.GenerateSubscriptionCacheKey(id)
	if err := uc.cache.Set(cacheKey, subscription, 24*time.Hour); err != nil {
		log.Printf("Warning: Failed to update subscription in cache: %v", err)
	}

	// Send renewal confirmation email
	if err := uc.emailService.SendRenewalEmail(subscription.CustomerEmail); err != nil {
		log.Printf("Warning: Failed to send renewal email: %v", err)
	}

	return nil
}

func (uc *subscriptionUseCase) ListSubscriptions(customerEmail string, page, limit int32) ([]*domain.Subscription, int32, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10 // Default limit
	}

	return uc.repo.List(customerEmail, page, limit)
}
