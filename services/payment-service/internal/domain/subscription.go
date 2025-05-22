package domain

import (
	"time"
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "ACTIVE"
	SubscriptionStatusCancelled SubscriptionStatus = "CANCELLED"
	SubscriptionStatusExpired   SubscriptionStatus = "EXPIRED"
)

type Subscription struct {
	ID            string
	CustomerEmail string
	PlanName      string
	Price         float64
	Currency      string
	Status        SubscriptionStatus
	StartDate     time.Time
	EndDate       time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type SubscriptionRepository interface {
	Create(subscription *Subscription) error
	GetByID(id string) (*Subscription, error)
	Update(subscription *Subscription) error
	Delete(id string) error
	List(customerEmail string, page, limit int32) ([]*Subscription, int32, error)
}

type SubscriptionCache interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (interface{}, error)
	Delete(key string) error
}

type SubscriptionUseCase interface {
	CreateSubscription(subscription *Subscription) error
	GetSubscription(id string) (*Subscription, error)
	CancelSubscription(id string) error
	RenewSubscription(id string) error
	ListSubscriptions(customerEmail string, page, limit int32) ([]*Subscription, int32, error)
}
