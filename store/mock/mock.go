package store

import (
	"context"

	"github.com/Ruthvik10/targeting-engine/model"
)

type deliveryStoreMock struct {
}

func NewDeliveryStoreMock() *deliveryStoreMock {
	return &deliveryStoreMock{}
}

func (store *deliveryStoreMock) GetCampaigns(ctx context.Context, in *model.Delivery) ([]*model.Campaign, error) {
	return GetCampaignsMock(ctx, in)
}

// GetCampaignsMock will be mocked in the handlers.
var GetCampaignsMock = func(ctx context.Context, in *model.Delivery) ([]*model.Campaign, error) {
	return nil, nil
}
