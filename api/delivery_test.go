package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockCache "github.com/Ruthvik10/targeting-engine/cache/mock"
	"github.com/Ruthvik10/targeting-engine/model"
	mockStore "github.com/Ruthvik10/targeting-engine/store/mock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Test cases
func TestDeliverCampaign(t *testing.T) {

	t.Run("should return 400 when app param is missing", func(t *testing.T) {
		// Initialize mocks
		store := mockStore.NewDeliveryStore()
		cache := mockCache.NewDeliveryCache()

		handler := NewDeliveryHandler(store, cache, 5*time.Second)
		req, _ := http.NewRequest(http.MethodGet, "/v1/delivery?country=US&os=Android", nil)
		rec := httptest.NewRecorder()

		handler.DeliverCampaign(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"error": "missing app param"}`, rec.Body.String())
	})

	t.Run("should return 400 when country param is missing", func(t *testing.T) {
		// Initialize mocks
		store := mockStore.NewDeliveryStore()
		cache := mockCache.NewDeliveryCache()

		handler := NewDeliveryHandler(store, cache, 5*time.Second)
		req, _ := http.NewRequest(http.MethodGet, "/v1/delivery?app=com.spotify&os=Android", nil)
		rec := httptest.NewRecorder()

		handler.DeliverCampaign(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"error": "missing country param"}`, rec.Body.String())
	})

	t.Run("should return 400 when os param is missing", func(t *testing.T) {
		// Initialize mocks
		store := mockStore.NewDeliveryStore()
		cache := mockCache.NewDeliveryCache()

		handler := NewDeliveryHandler(store, cache, 5*time.Second)
		req, _ := http.NewRequest(http.MethodGet, "/v1/delivery?app=com.spotify&country=US", nil)
		rec := httptest.NewRecorder()

		handler.DeliverCampaign(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"error": "missing os param"}`, rec.Body.String())
	})

	t.Run("should return campaigns from cache on cache hit", func(t *testing.T) {
		// Initialize mocks
		store := mockStore.NewDeliveryStore()
		cache := mockCache.NewDeliveryCache()

		handler := NewDeliveryHandler(store, cache, 5*time.Second)
		campaign1ID := primitive.NewObjectID()
		campaign2ID := primitive.NewObjectID()
		// Mock cache hit
		cache.GetCampaignsFunc = func(ctx context.Context, key string) (string, error) {
			campaign1 := &model.Campaign{
				ID:    campaign1ID,
				Name:  "Spotify Campaign",
				Image: "https://somelink1",
				CTA:   "Download",
				Targeting: model.Targeting{
					IncludeOS:      []string{"Android"},
					IncludeCountry: []string{"US"},
					IncludeApp:     []string{"com.spotify"},
					ExcludeApp:     []string{},
					ExcludeOS:      []string{},
					ExcludeCountry: []string{},
				},
			}

			campaign2 := &model.Campaign{
				ID:    campaign2ID,
				Name:  "Spotify Campaign V2",
				Image: "https://somelink2",
				CTA:   "Install",
				Targeting: model.Targeting{
					IncludeOS:      []string{},
					IncludeCountry: []string{},
					IncludeApp:     []string{},
					ExcludeApp:     []string{"com.gametion.ludokinggame"},
					ExcludeOS:      []string{"iOS"},
					ExcludeCountry: []string{"UK"},
				},
			}
			campaigns := []*model.Campaign{campaign1, campaign2}
			toReturn, _ := json.Marshal(campaigns)
			return string(toReturn), nil
		}

		req, _ := http.NewRequest(http.MethodGet, "/v1/delivery?app=com.spotify&os=android&country=us", nil)
		rec := httptest.NewRecorder()

		handler.DeliverCampaign(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		expectedResponse := []*DeliveryResponse{
			{
				CID:   campaign1ID.Hex(),
				Image: "https://somelink1",
				CTA:   "Download",
			},
			{
				CID:   campaign2ID.Hex(),
				Image: "https://somelink2",
				CTA:   "Install",
			},
		}
		expectedJSON, _ := json.Marshal(expectedResponse)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("should fetch campaigns from DB on cache miss and cache the result", func(t *testing.T) {
		// Initialize mocks
		store := mockStore.NewDeliveryStore()
		cache := mockCache.NewDeliveryCache()

		handler := NewDeliveryHandler(store, cache, 5*time.Second)

		// Mock cache miss and DB hit

		cache.GetCampaignsFunc = func(ctx context.Context, key string) (string, error) {
			return "", nil
		}

		cache.SetCampaignFunc = func(ctx context.Context, key, value string, exp time.Duration) error {
			return nil
		}

		campaign1ID := primitive.NewObjectID()
		campaign2ID := primitive.NewObjectID()

		store.GetCampaignsFunc = func(ctx context.Context, in *model.Delivery) ([]*model.Campaign, error) {
			campaign1 := &model.Campaign{
				ID:    campaign1ID,
				Name:  "Spotify Campaign",
				Image: "https://somelink1",
				CTA:   "Download",
				Targeting: model.Targeting{
					IncludeOS:      []string{"Android"},
					IncludeCountry: []string{"US"},
					IncludeApp:     []string{"com.spotify"},
					ExcludeApp:     []string{},
					ExcludeOS:      []string{},
					ExcludeCountry: []string{},
				},
			}

			campaign2 := &model.Campaign{
				ID:    campaign2ID,
				Name:  "Spotify Campaign V2",
				Image: "https://somelink2",
				CTA:   "Install",
				Targeting: model.Targeting{
					IncludeOS:      []string{},
					IncludeCountry: []string{},
					IncludeApp:     []string{},
					ExcludeApp:     []string{"com.gametion.ludokinggame"},
					ExcludeOS:      []string{"iOS"},
					ExcludeCountry: []string{"UK"},
				},
			}
			dbCampaigns := []*model.Campaign{campaign1, campaign2}
			return dbCampaigns, nil
		}

		req, _ := http.NewRequest(http.MethodGet, "/v1/delivery?app=com.spotify&os=android&country=us", nil)
		rec := httptest.NewRecorder()

		handler.DeliverCampaign(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		expectedResponse := []*DeliveryResponse{
			{
				CID:   campaign1ID.Hex(),
				Image: "https://somelink1",
				CTA:   "Download",
			},
			{
				CID:   campaign2ID.Hex(),
				Image: "https://somelink2",
				CTA:   "Install",
			},
		}
		expectedJSON, _ := json.Marshal(expectedResponse)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())

	})

	t.Run("should query the database if something goes wrong while unmarshalling data from cache", func(t *testing.T) {
		store := mockStore.NewDeliveryStore()
		cache := mockCache.NewDeliveryCache()

		handler := NewDeliveryHandler(store, cache, 5*time.Second)
		campaign1ID := primitive.NewObjectID()
		campaign2ID := primitive.NewObjectID()

		// Mock invalid data present in the cache.
		// NOTE: This scenario should not happen under normal circumstances.
		cache.GetCampaignsFunc = func(ctx context.Context, key string) (string, error) {

			toReturn := `{"invalidField1": "val1"}`
			return toReturn, nil
		}

		// Mock database call to get campaigns
		store.GetCampaignsFunc = func(ctx context.Context, in *model.Delivery) ([]*model.Campaign, error) {
			campaign1 := &model.Campaign{
				ID:    campaign1ID,
				Name:  "Spotify Campaign",
				Image: "https://somelink1",
				CTA:   "Download",
				Targeting: model.Targeting{
					IncludeOS:      []string{"Android"},
					IncludeCountry: []string{"US"},
					IncludeApp:     []string{"com.spotify"},
					ExcludeApp:     []string{},
					ExcludeOS:      []string{},
					ExcludeCountry: []string{},
				},
			}

			campaign2 := &model.Campaign{
				ID:    campaign2ID,
				Name:  "Spotify Campaign V2",
				Image: "https://somelink2",
				CTA:   "Install",
				Targeting: model.Targeting{
					IncludeOS:      []string{},
					IncludeCountry: []string{},
					IncludeApp:     []string{},
					ExcludeApp:     []string{"com.gametion.ludokinggame"},
					ExcludeOS:      []string{"iOS"},
					ExcludeCountry: []string{"UK"},
				},
			}
			dbCampaigns := []*model.Campaign{campaign1, campaign2}
			return dbCampaigns, nil
		}

		req, _ := http.NewRequest(http.MethodGet, "/v1/delivery?app=com.spotify&os=android&country=us", nil)
		rec := httptest.NewRecorder()

		handler.DeliverCampaign(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		expectedResponse := []*DeliveryResponse{
			{
				CID:   campaign1ID.Hex(),
				Image: "https://somelink1",
				CTA:   "Download",
			},
			{
				CID:   campaign2ID.Hex(),
				Image: "https://somelink2",
				CTA:   "Install",
			},
		}
		expectedJSON, _ := json.Marshal(expectedResponse)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("should return 204 when no campaigns are found", func(t *testing.T) {

		// Initialize mocks
		store := mockStore.NewDeliveryStore()
		cache := mockCache.NewDeliveryCache()

		handler := NewDeliveryHandler(store, cache, 5*time.Second)

		// Mock cache miss and DB hit returning no campaigns

		cache.GetCampaignsFunc = func(ctx context.Context, key string) (string, error) {
			return "", nil
		}
		store.GetCampaignsFunc = func(ctx context.Context, in *model.Delivery) ([]*model.Campaign, error) {
			return []*model.Campaign{}, nil
		}
		req, _ := http.NewRequest(http.MethodGet, "/v1/delivery?app=com.spotify&os=android&country=uk", nil)
		rec := httptest.NewRecorder()

		handler.DeliverCampaign(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("should return 500 if DB query fails", func(t *testing.T) {

		// Initialize mocks
		store := mockStore.NewDeliveryStore()
		cache := mockCache.NewDeliveryCache()

		handler := NewDeliveryHandler(store, cache, 5*time.Second)

		// Mock cache miss and DB query failure
		cache.GetCampaignsFunc = func(ctx context.Context, key string) (string, error) {
			return "", nil
		}
		store.GetCampaignsFunc = func(ctx context.Context, in *model.Delivery) ([]*model.Campaign, error) {
			return []*model.Campaign{}, errors.New("error")
		}

		req, _ := http.NewRequest(http.MethodGet, "/v1/delivery?app=com.spotify&os=android&country=uk", nil)
		rec := httptest.NewRecorder()

		handler.DeliverCampaign(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t, `{"error": "Something went wrong while fetching the campaigns!"}`, rec.Body.String())
	})
}
