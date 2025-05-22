package cache

import (
	"context"
	"encoding/json"
	"time"

	"payment-service/internal/domain"

	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) domain.PaymentCache {
	return &redisCache{client: client}
}

func (c *redisCache) Set(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, expiration).Err()
}

func (c *redisCache) Get(key string) (interface{}, error) {
	ctx := context.Background()
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var payment domain.Payment
	if err := json.Unmarshal(data, &payment); err != nil {
		return nil, err
	}
	return &payment, nil
}

func (c *redisCache) Delete(key string) error {
	ctx := context.Background()
	return c.client.Del(ctx, key).Err()
}

// Helper method to generate cache key for payment
func GeneratePaymentCacheKey(paymentID string) string {
	return "payment:" + paymentID
}
