package cache

import (
	"context"
	"encoding/json"
	"time"

	"payment-service/internal/domain"

	"github.com/redis/go-redis/v9"
)

type subscriptionCache struct {
	client *redis.Client
}

func NewSubscriptionCache(client *redis.Client) domain.SubscriptionCache {
	return &subscriptionCache{client: client}
}

func (c *subscriptionCache) Set(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, expiration).Err()
}

func (c *subscriptionCache) Get(key string) (interface{}, error) {
	ctx := context.Background()
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var subscription domain.Subscription
	if err := json.Unmarshal(data, &subscription); err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (c *subscriptionCache) Delete(key string) error {
	ctx := context.Background()
	return c.client.Del(ctx, key).Err()
}

// Helper function to generate subscription cache key
func GenerateSubscriptionCacheKey(subscriptionID string) string {
	return "subscription:" + subscriptionID
}
