package mock

import (
	"context"
	"time"
)

type deliveryCache struct {
	CountKeysCalled       bool
	DeleteCampaignCalled  bool
	CountKeysReturnKeys   []string
	CountKeysReturnError  error
	DeleteCampaignResults map[string]error
	CountKeysFunc         func(ctx context.Context, key string) ([]string, error)
	DeleteCampaignFunc    func(ctx context.Context, key string) error
	GetCampaignsFunc      func(ctx context.Context, key string) (string, error)
	SetCampaignFunc       func(ctx context.Context, key, value string, exp time.Duration) error
}

func NewDeliveryCache() *deliveryCache {
	return &deliveryCache{}
}

func (c *deliveryCache) GetCampaigns(ctx context.Context, key string) (string, error) {
	return c.GetCampaignsFunc(ctx, key)
}

func (c *deliveryCache) SetCampaign(ctx context.Context, key, value string, exp time.Duration) error {
	return c.SetCampaignFunc(ctx, key, value, exp)
}

func (c *deliveryCache) CountKeys(ctx context.Context, key string) ([]string, error) {
	return c.CountKeysFunc(ctx, key)
}

func (c *deliveryCache) DeleteCampaign(ctx context.Context, key string) error {
	return c.DeleteCampaignFunc(ctx, key)
}
