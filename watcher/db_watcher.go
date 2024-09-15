package watcher

import (
	"context"
	"fmt"
	"log"

	"github.com/Ruthvik10/targeting-engine/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type store interface {
	GetCampaigns(ctx context.Context, in *model.Delivery) ([]*model.Campaign, error)
	WatchCampaign(ctx context.Context, out chan<- bson.M)
}

type cache interface {
	CountKeys(ctx context.Context, key string) ([]string, error)
	DeleteCampaign(ctx context.Context, key string) error
}

type DBWatcher struct {
	store store
	cache cache
}

func NewDBWatcher(store store, cache cache) *DBWatcher {
	return &DBWatcher{store: store, cache: cache}
}

func (w *DBWatcher) WatchCampaign(ctx context.Context) {
	changedDocCh := make(chan bson.M)
	go w.store.WatchCampaign(ctx, changedDocCh)
	go func(chan bson.M) {
		for {
			changedDoc := <-changedDocCh
			handleCampaignUpdate(ctx, changedDoc, w.cache)

		}
	}(changedDocCh)
}

func handleCampaignUpdate(ctx context.Context, event bson.M, cache cache) {
	fullDocument := event["fullDocument"].(bson.M)
	campaignID := fullDocument["_id"].(primitive.ObjectID)
	status := fullDocument["status"].(string)

	if status == "INACTIVE" {
		// Fetch all delivery cache keys (app:country:os) associated with this campaign
		campaignKey := fmt.Sprintf("campaign:%s", campaignID)
		deliveryCacheKeys, err := cache.CountKeys(ctx, campaignKey)
		if err != nil {
			log.Printf("Error fetching cache keys for campaign: %v", err)
			return
		}

		// Invalidate each cache key
		for _, cacheKey := range deliveryCacheKeys {
			err := cache.DeleteCampaign(ctx, cacheKey)
			if err != nil {
				log.Printf("Failed to delete cache key: %v", cacheKey)
			} else {
				log.Printf("Deleted cache key: %v", cacheKey)
			}
		}

		cache.DeleteCampaign(ctx, campaignKey)
	}
}
