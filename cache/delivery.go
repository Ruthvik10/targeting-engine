package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Ruthvik10/targeting-engine/model"
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
	campaigns := make([]*model.Campaign, 0)
	err := json.Unmarshal([]byte(value), &campaigns)
	if err != nil {
		log.Printf("Error unmarshalling campaigns: %v\n", err)
	}

	// We maintain two caches.
	// One is a cache with key "campaign:campaignID" and value which is a set [app1:country1:os2, app1:country2:os1].
	// Second is the a cache with key "app:country:os" and value which is "[campaign1, campaign2]"
	// The first cache is useful during the deletion process to identify all the deliveries that a campaign serves and delete it, for ex, when a campaign becomes inactive.
	for _, campaign := range campaigns {
		campaignKey := fmt.Sprintf("campaign:%s", campaign.ID)
		// Add the cache key to the set for this campaign
		c.client.SAdd(ctx, campaignKey, key)
	}
	err = c.client.Set(ctx, key, value, exp).Err()
	if err != nil {
		log.Printf("Error setting the cache value for the key [%s]: %v", key, err)
	}
	return err
}

func (c *Delivery) CountKeys(ctx context.Context, key string) ([]string, error) {
	keys, err := c.client.SMembers(ctx, key).Result()
	if err != nil {
		log.Printf("Error getting count of the keys for the campaign key (%s): %v\n", key, err)
		return nil, err
	}
	return keys, nil
}

func (c *Delivery) DeleteCampaign(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		log.Printf("Error deleting the key (%s): %v\n", key, err)
	}

	return err
}
