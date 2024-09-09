package cache

import (
	"context"
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
	return c.client.Get(ctx, key).Result()
}

func (c *Delivery) SetCampaign(ctx context.Context, key, value string, exp time.Duration) error {
	return c.client.Set(ctx, key, value, exp).Err()
}
