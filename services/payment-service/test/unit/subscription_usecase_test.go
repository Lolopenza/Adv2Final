package unit

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"payment-service/internal/domain"
	"payment-service/internal/usecase"
)

// Mock implementations for subscription testing
type MockSubscriptionRepository struct {
	mock.Mock
}

func (m *MockSubscriptionRepository) Create(subscription *domain.Subscription) error {
	args := m.Called(subscription)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) GetByID(id string) (*domain.Subscription, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) Update(subscription *domain.Subscription) error {
	args := m.Called(subscription)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) List(customerEmail string, page, limit int32) ([]*domain.Subscription, int32, error) {
	args := m.Called(customerEmail, page, limit)
	return args.Get(0).([]*domain.Subscription), args.Get(1).(int32), args.Error(2)
}

type MockSubscriptionCache struct {
	mock.Mock
}

func (m *MockSubscriptionCache) Set(key string, value interface{}, expiration time.Duration) error {
	args := m.Called(key, value, expiration)
	return args.Error(0)
}

func (m *MockSubscriptionCache) Get(key string) (interface{}, error) {
	args := m.Called(key)
	return args.Get(0), args.Error(1)
}

func (m *MockSubscriptionCache) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

// Test cases
func TestCreateSubscription(t *testing.T) {
	// Setup
	mockRepo := new(MockSubscriptionRepository)
	mockCache := new(MockSubscriptionCache)
	mockEventPublisher := new(MockPaymentEventPublisher)
	mockEmail := new(MockEmailService)

	useCase := usecase.NewSubscriptionUseCase(mockRepo, mockCache, mockEventPublisher, mockEmail)

	// Test data
	subscription := &domain.Subscription{
		CustomerEmail: "test@example.com",
		PlanName:      "Premium",
		Price:         19.99,
		Currency:      "USD",
	}

	// Expectations
	mockRepo.On("Create", mock.AnythingOfType("*domain.Subscription")).Return(nil)
	mockCache.On("Set", mock.AnythingOfType("string"), mock.AnythingOfType("*domain.Subscription"), mock.AnythingOfType("time.Duration")).Return(nil)

	// Execute
	err := useCase.CreateSubscription(subscription)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, subscription.ID)
	assert.Equal(t, domain.SubscriptionStatusActive, subscription.Status)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestGetSubscription(t *testing.T) {
	// Setup
	mockRepo := new(MockSubscriptionRepository)
	mockCache := new(MockSubscriptionCache)
	mockEventPublisher := new(MockPaymentEventPublisher)
	mockEmail := new(MockEmailService)

	useCase := usecase.NewSubscriptionUseCase(mockRepo, mockCache, mockEventPublisher, mockEmail)

	// Test data
	subscriptionID := uuid.New().String()
	subscription := &domain.Subscription{
		ID:            subscriptionID,
		CustomerEmail: "test@example.com",
		PlanName:      "Premium",
		Price:         19.99,
		Currency:      "USD",
		Status:        domain.SubscriptionStatusActive,
		StartDate:     time.Now(),
		EndDate:       time.Now().AddDate(0, 1, 0),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Test Case 1: Cache hit
	mockCache.On("Get", mock.AnythingOfType("string")).Return(subscription, nil).Once()

	// Execute
	result, err := useCase.GetSubscription(subscriptionID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, subscription, result)

	// Test Case 2: Cache miss, DB hit
	mockCache.On("Get", mock.AnythingOfType("string")).Return(nil, errors.New("cache miss")).Once()
	mockRepo.On("GetByID", subscriptionID).Return(subscription, nil).Once()
	mockCache.On("Set", mock.AnythingOfType("string"), mock.AnythingOfType("*domain.Subscription"), mock.AnythingOfType("time.Duration")).Return(nil).Once()

	// Execute
	result, err = useCase.GetSubscription(subscriptionID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, subscription, result)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestCancelSubscription(t *testing.T) {
	// Setup
	mockRepo := new(MockSubscriptionRepository)
	mockCache := new(MockSubscriptionCache)
	mockEventPublisher := new(MockPaymentEventPublisher)
	mockEmail := new(MockEmailService)

	useCase := usecase.NewSubscriptionUseCase(mockRepo, mockCache, mockEventPublisher, mockEmail)

	// Test data
	subscriptionID := uuid.New().String()
	subscription := &domain.Subscription{
		ID:            subscriptionID,
		CustomerEmail: "test@example.com",
		PlanName:      "Premium",
		Price:         19.99,
		Currency:      "USD",
		Status:        domain.SubscriptionStatusActive,
		StartDate:     time.Now(),
		EndDate:       time.Now().AddDate(0, 1, 0),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Expectations
	mockCache.On("Get", mock.AnythingOfType("string")).Return(subscription, nil).Once()
	mockRepo.On("Update", mock.AnythingOfType("*domain.Subscription")).Run(func(args mock.Arguments) {
		sub := args.Get(0).(*domain.Subscription)
		assert.Equal(t, domain.SubscriptionStatusCancelled, sub.Status)
	}).Return(nil).Once()
	mockCache.On("Set", mock.AnythingOfType("string"), mock.AnythingOfType("*domain.Subscription"), mock.AnythingOfType("time.Duration")).Return(nil).Once()
	mockEmail.On("SendCancellationEmail", subscription.CustomerEmail).Return(nil).Once()

	// Execute
	err := useCase.CancelSubscription(subscriptionID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
}

func TestRenewSubscription(t *testing.T) {
	// Setup
	mockRepo := new(MockSubscriptionRepository)
	mockCache := new(MockSubscriptionCache)
	mockEventPublisher := new(MockPaymentEventPublisher)
	mockEmail := new(MockEmailService)

	useCase := usecase.NewSubscriptionUseCase(mockRepo, mockCache, mockEventPublisher, mockEmail)

	// Test data
	subscriptionID := uuid.New().String()
	oldEndDate := time.Now()
	subscription := &domain.Subscription{
		ID:            subscriptionID,
		CustomerEmail: "test@example.com",
		PlanName:      "Premium",
		Price:         19.99,
		Currency:      "USD",
		Status:        domain.SubscriptionStatusActive,
		StartDate:     time.Now().Add(-30 * 24 * time.Hour), // 30 days ago
		EndDate:       oldEndDate,
		CreatedAt:     time.Now().Add(-30 * 24 * time.Hour),
		UpdatedAt:     time.Now().Add(-30 * 24 * time.Hour),
	}

	// Expectations
	mockCache.On("Get", mock.AnythingOfType("string")).Return(subscription, nil).Once()
	mockRepo.On("Update", mock.AnythingOfType("*domain.Subscription")).Run(func(args mock.Arguments) {
		sub := args.Get(0).(*domain.Subscription)
		assert.True(t, sub.EndDate.After(oldEndDate), "End date should be extended")
	}).Return(nil).Once()
	mockCache.On("Set", mock.AnythingOfType("string"), mock.AnythingOfType("*domain.Subscription"), mock.AnythingOfType("time.Duration")).Return(nil).Once()
	mockEmail.On("SendRenewalEmail", subscription.CustomerEmail).Return(nil).Once()

	// Execute
	err := useCase.RenewSubscription(subscriptionID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
}

func TestListSubscriptions(t *testing.T) {
	// Setup
	mockRepo := new(MockSubscriptionRepository)
	mockCache := new(MockSubscriptionCache)
	mockEventPublisher := new(MockPaymentEventPublisher)
	mockEmail := new(MockEmailService)

	useCase := usecase.NewSubscriptionUseCase(mockRepo, mockCache, mockEventPublisher, mockEmail)

	// Test data
	customerEmail := "test@example.com"
	subscriptions := []*domain.Subscription{
		{
			ID:            uuid.New().String(),
			CustomerEmail: customerEmail,
			PlanName:      "Premium",
			Status:        domain.SubscriptionStatusActive,
		},
		{
			ID:            uuid.New().String(),
			CustomerEmail: customerEmail,
			PlanName:      "Basic",
			Status:        domain.SubscriptionStatusCancelled,
		},
	}
	var total int32 = 2

	// Expectations
	mockRepo.On("List", customerEmail, int32(1), int32(10)).Return(subscriptions, total, nil).Once()

	// Execute
	result, resultTotal, err := useCase.ListSubscriptions(customerEmail, 1, 10)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, subscriptions, result)
	assert.Equal(t, total, resultTotal)
	mockRepo.AssertExpectations(t)
}
