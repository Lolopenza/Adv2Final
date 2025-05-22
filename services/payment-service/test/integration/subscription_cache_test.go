package integration

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"payment-service/internal/domain"
	"payment-service/internal/repository/cache"
)

func TestSubscriptionCache(t *testing.T) {
	// Start miniredis server
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer s.Close()

	// Create Redis client connected to miniredis
	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	// Initialize cache repository
	subscriptionCache := cache.NewSubscriptionCache(client)

	// Create test subscription
	subscription := &domain.Subscription{
		ID:            uuid.New().String(),
		CustomerEmail: "test@example.com",
		PlanName:      "Premium",
		Price:         19.99,
		Currency:      "USD",
		Status:        domain.SubscriptionStatusActive,
		StartDate:     time.Now(),
		EndDate:       time.Now().AddDate(0, 1, 0), // 1 month
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Test 1: Set and Get
	t.Run("Set and Get", func(t *testing.T) {
		cacheKey := cache.GenerateSubscriptionCacheKey(subscription.ID)

		// Set subscription in cache
		err := subscriptionCache.Set(cacheKey, subscription, 1*time.Hour)
		require.NoError(t, err)

		// Get subscription from cache
		cachedSub, err := subscriptionCache.Get(cacheKey)
		require.NoError(t, err)
		require.NotNil(t, cachedSub)

		// Verify retrieved subscription
		retrievedSub := cachedSub.(*domain.Subscription)
		assert.Equal(t, subscription.ID, retrievedSub.ID)
		assert.Equal(t, subscription.CustomerEmail, retrievedSub.CustomerEmail)
		assert.Equal(t, subscription.PlanName, retrievedSub.PlanName)
		assert.Equal(t, subscription.Status, retrievedSub.Status)
	})

	// Test 2: Cache miss
	t.Run("Cache Miss", func(t *testing.T) {
		missingKey := cache.GenerateSubscriptionCacheKey(uuid.New().String())

		// Try to get non-existent subscription
		cachedSub, err := subscriptionCache.Get(missingKey)
		assert.Nil(t, cachedSub)
		// Nil значение указывает на отсутствие в кэше согласно реализации
		assert.Nil(t, err)
	})

	// Test 3: Delete (cache invalidation)
	t.Run("Delete", func(t *testing.T) {
		cacheKey := cache.GenerateSubscriptionCacheKey(subscription.ID)

		// Set subscription in cache
		err := subscriptionCache.Set(cacheKey, subscription, 1*time.Hour)
		require.NoError(t, err)

		// Delete from cache
		err = subscriptionCache.Delete(cacheKey)
		require.NoError(t, err)

		// Verify it's gone
		cachedSub, err := subscriptionCache.Get(cacheKey)
		assert.Nil(t, cachedSub)
		assert.Nil(t, err)
	})

	// Test 4: Update (cache invalidation and rewrite)
	t.Run("Update", func(t *testing.T) {
		cacheKey := cache.GenerateSubscriptionCacheKey(subscription.ID)

		// Set original subscription
		err := subscriptionCache.Set(cacheKey, subscription, 1*time.Hour)
		require.NoError(t, err)

		// Update subscription
		updatedSub := *subscription
		updatedSub.Status = domain.SubscriptionStatusCancelled
		updatedSub.UpdatedAt = time.Now()

		// Set updated subscription
		err = subscriptionCache.Set(cacheKey, &updatedSub, 1*time.Hour)
		require.NoError(t, err)

		// Get updated subscription
		cachedSub, err := subscriptionCache.Get(cacheKey)
		require.NoError(t, err)

		// Verify updated status
		retrievedSub := cachedSub.(*domain.Subscription)
		assert.Equal(t, domain.SubscriptionStatusCancelled, retrievedSub.Status)
	})

	// Test 5: TTL Expiration
	t.Run("TTL Expiration", func(t *testing.T) {
		cacheKey := cache.GenerateSubscriptionCacheKey(subscription.ID)

		// Set with short TTL
		err := subscriptionCache.Set(cacheKey, subscription, 1*time.Second)
		require.NoError(t, err)

		// Fast forward time in miniredis
		s.FastForward(2 * time.Second)

		// Check that key is expired
		cachedSub, err := subscriptionCache.Get(cacheKey)
		assert.Nil(t, cachedSub)
		assert.Nil(t, err) // Согласно реализации cache, когда ключ не найден, возвращается nil, nil
	})
}
