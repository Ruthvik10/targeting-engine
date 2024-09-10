package cache

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Delivery struct {
	client *redis.Client
}

func NewDelivery(client *redis.Client) *Delivery {
	return &Delivery{client: client}
}

func (c *Delivery) GetCampaigns(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		log.Printf("Error getting the cached value for the key [%s]: %v", key, err)
	}
	return val, err
}

func (c *Delivery) SetCampaign(ctx context.Context, key, value string, exp time.Duration) error {
	err := c.client.Set(ctx, key, value, exp).Err()
	if err != nil {
		log.Printf("Error setting the cache value for the key [%s]: %v", key, err)
	}
	return err
}
