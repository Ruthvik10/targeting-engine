package mock

import (
	"context"

	"github.com/Ruthvik10/targeting-engine/model"
	"go.mongodb.org/mongo-driver/bson"
)

type deliveryStore struct {
	WatchCampaignFunc func(ctx context.Context, out chan<- bson.M)
	GetCampaignsFunc  func(ctx context.Context, in *model.Delivery) ([]*model.Campaign, error)
}

func NewDeliveryStore() *deliveryStore {
	return &deliveryStore{}
}

func (store *deliveryStore) GetCampaigns(ctx context.Context, in *model.Delivery) ([]*model.Campaign, error) {
	return store.GetCampaignsFunc(ctx, in)
}

func (store *deliveryStore) WatchCampaign(ctx context.Context, out chan<- bson.M) {
	store.WatchCampaignFunc(ctx, out)
}
