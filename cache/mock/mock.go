package mock

import (
	"context"
	"time"
)

type deliveryCacheMock struct {
}

func NewDeliveryCacheMock() *deliveryCacheMock {
	return &deliveryCacheMock{}
}

func (c *deliveryCacheMock) GetCampaigns(ctx context.Context, key string) (string, error) {
	return GetCampaignsMock(ctx, key)
}

func (c *deliveryCacheMock) SetCampaign(ctx context.Context, key, value string, exp time.Duration) error {
	return SetCampaignMock(ctx, key, value, exp)
}

// GetCampaignsMock will be mocked in the api handler.
var GetCampaignsMock = func(ctx context.Context, key string) (string, error) {
	return "", nil
}

// SetCampaignMock will be mocked in the api handler.
var SetCampaignMock = func(ctx context.Context, key, value string, exp time.Duration) error {
	return nil
}
