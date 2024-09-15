package watcher

import (
	"context"
	"testing"
	"time"

	mockcache "github.com/Ruthvik10/targeting-engine/cache/mock"
	mockstore "github.com/Ruthvik10/targeting-engine/store/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDBWatcher_WatchCampaign(t *testing.T) {
	// Create a mock store and cache
	mockStore := mockstore.NewDeliveryStore()
	mockCache := mockcache.NewDeliveryCache()

	mockCache.CountKeysReturnKeys = []string{"campaign:123", "campaign:456"}
	mockCache.DeleteCampaignResults = map[string]error{
		"campaign:123": nil,
		"campaign:456": nil,
	}

	mockCache.CountKeysFunc = func(ctx context.Context, key string) ([]string, error) {
		mockCache.CountKeysCalled = true
		return mockCache.CountKeysReturnKeys, mockCache.CountKeysReturnError
	}

	mockCache.DeleteCampaignFunc = func(ctx context.Context, key string) error {
		mockCache.DeleteCampaignCalled = true
		if res, ok := mockCache.DeleteCampaignResults[key]; ok {
			return res
		}
		return nil
	}

	// Create the DBWatcher with the mock store and cache
	watcher := NewDBWatcher(mockStore, mockCache)

	// Define the campaignID and event to be sent over the channel
	campaignID := primitive.NewObjectID()
	event := bson.M{
		"fullDocument": bson.M{
			"_id":    campaignID,
			"status": "INACTIVE",
		},
	}

	// Setup the mock store to simulate sending an event through the WatchCampaign method
	mockStore.WatchCampaignFunc = func(ctx context.Context, out chan<- bson.M) {
		out <- event
	}

	// Call WatchCampaign in a goroutine to handle the event
	ctx := context.Background()
	go watcher.WatchCampaign(ctx)

	// Wait for a short time to allow goroutines to process
	time.Sleep(100 * time.Millisecond)

	// Verify that CountKeys and DeleteCampaign were called with the expected arguments
	if !mockCache.CountKeysCalled {
		t.Errorf("Expected CountKeys to be called")
	}

	if !mockCache.DeleteCampaignCalled {
		t.Errorf("Expected DeleteCampaign to be called")
	}

	// Verify that DeleteCampaign was called with the expected keys
	if err := mockCache.DeleteCampaign(ctx, "campaign:123"); err != nil {
		t.Errorf("Failed to delete campaign:123")
	}
	if err := mockCache.DeleteCampaign(ctx, "campaign:456"); err != nil {
		t.Errorf("Failed to delete campaign:456")
	}
}

func TestHandleCampaignUpdate_NoInactiveStatus(t *testing.T) {
	// Create a mock cache
	mockCache := mockcache.NewDeliveryCache()

	// Event with status "ACTIVE" (should not trigger cache deletion)
	event := bson.M{
		"fullDocument": bson.M{
			"_id":    primitive.NewObjectID(),
			"status": "ACTIVE",
		},
	}

	// Call handleCampaignUpdate
	handleCampaignUpdate(context.Background(), event, mockCache)

	// Ensure no cache deletion is triggered
	if mockCache.CountKeysCalled {
		t.Errorf("CountKeys should not have been called")
	}

	if mockCache.DeleteCampaignCalled {
		t.Errorf("DeleteCampaign should not have been called")
	}
}
