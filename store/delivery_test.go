package store

import (
	"context"
	"fmt"
	"testing"

	"github.com/Ruthvik10/targeting-engine/model"
	"github.com/stretchr/testify/assert"
)

func createCampaigns(t *testing.T, campaigns []interface{}) {

	// Insert multiple documents
	insertManyResult, err := deliveryStore.collection.InsertMany(context.Background(), campaigns)
	assert.NoError(t, err)
	assert.NotEmpty(t, insertManyResult)
}

func TestGetCampaigns(t *testing.T) {
	newCampaigns := []interface{}{
		model.Campaign{
			Name:   "TuneUp - Your Music Companion",
			Image:  "https://somelink/tuneup.jpg",
			CTA:    "Listen Now",
			Status: "ACTIVE",
			Targeting: model.Targeting{
				IncludeApp:     []string{"com.tuneup.music"},
				IncludeCountry: []string{"UK", "Australia"},
				IncludeOS:      []string{"iOS", "Android"},
				ExcludeOS:      []string{},
				ExcludeApp:     []string{},
				ExcludeCountry: []string{},
			},
		},
		model.Campaign{
			Name:   "LinguaMaster - Language Learning",
			Image:  "https://somelink/linguamaster.jpg",
			CTA:    "Start Learning",
			Status: "ACTIVE",
			Targeting: model.Targeting{
				IncludeApp:     []string{"com.linguamaster.learn"},
				ExcludeCountry: []string{"France"},
				IncludeOS:      []string{"Android"},
				ExcludeApp:     []string{},
				ExcludeOS:      []string{},
				IncludeCountry: []string{},
			},
		},
		model.Campaign{
			Name:   "JetRun - Endless Runner Game",
			Image:  "https://somelink/jetrun.jpg",
			CTA:    "Play Now",
			Status: "ACTIVE",
			Targeting: model.Targeting{
				IncludeApp:     []string{"com.jetrun.game"},
				IncludeCountry: []string{"Germany", "Mexico"},
				ExcludeOS:      []string{"iOS"},
				IncludeOS:      []string{},
				ExcludeApp:     []string{},
				ExcludeCountry: []string{},
			},
		},
		model.Campaign{
			Name:   "GTA-IV",
			Image:  "https://somelink/gta.jpg",
			CTA:    "Play Now",
			Status: "INACTIVE",
			Targeting: model.Targeting{
				IncludeApp:     []string{"com.gtaiv.game"},
				IncludeCountry: []string{"US"},
				IncludeOS:      []string{"iOS"},
				ExcludeOS:      []string{},
				ExcludeApp:     []string{},
				ExcludeCountry: []string{},
			},
		},
		model.Campaign{
			Name:   "Counter Strike",
			Image:  "https://somelink/cs.jpg",
			CTA:    "Play Now",
			Status: "ACTIVE",
			Targeting: model.Targeting{
				ExcludeApp:     []string{"com.csg.game", "com.gtaiv.game"},
				IncludeCountry: []string{"US"},
				IncludeOS:      []string{"iOS"},
				ExcludeOS:      []string{},
				IncludeApp:     []string{},
				ExcludeCountry: []string{},
			},
		},
	}
	createCampaigns(t, newCampaigns)

	t.Run("Should not return any campaigns if the campaigns are inactive", func(t *testing.T) {

		query := model.Delivery{
			AppID:   "com.gtaiv.game",
			Country: "US",
			OS:      "iOS",
		}
		campaigns, err := deliveryStore.GetCampaigns(context.Background(), &query)
		assert.NoError(t, err)
		assert.Empty(t, campaigns)
	})
	t.Run("Should return the campaigns matching the query", func(t *testing.T) {

		query := model.Delivery{
			AppID:   "com.tuneup.music",
			Country: "UK",
			OS:      "Android",
		}
		campaigns, err := deliveryStore.GetCampaigns(context.Background(), &query)
		fmt.Println(campaigns)
		assert.NoError(t, err)
		assert.Equal(t, "TuneUp - Your Music Companion", campaigns[0].Name)
	})

	t.Run("Should not return any campaigns if the OS is in the exclusion list", func(t *testing.T) {

		query := model.Delivery{
			AppID:   "com.jetrun.game",
			Country: "Germany",
			OS:      "iOS",
		}
		campaigns, err := deliveryStore.GetCampaigns(context.Background(), &query)
		fmt.Println(campaigns)
		assert.NoError(t, err)
		assert.Empty(t, campaigns)
	})

	t.Run("Should not return any campaigns if the OS is not in the inclusion list", func(t *testing.T) {

		query := model.Delivery{
			AppID:   "com.linguamaster.learn",
			Country: "UK",
			OS:      "iOS",
		}
		campaigns, err := deliveryStore.GetCampaigns(context.Background(), &query)
		assert.NoError(t, err)
		assert.Empty(t, campaigns)
	})

	t.Run("Should not return any campaigns if the Country is in the exclusion list", func(t *testing.T) {

		query := model.Delivery{
			AppID:   "com.linguamaster.learn",
			Country: "France",
			OS:      "Android",
		}
		campaigns, err := deliveryStore.GetCampaigns(context.Background(), &query)
		fmt.Println(campaigns)
		assert.NoError(t, err)
		assert.Empty(t, campaigns)
	})

	t.Run("Should not return any campaigns if the Country is not in the inclusion list", func(t *testing.T) {

		query := model.Delivery{
			AppID:   "com.jetrun.game",
			Country: "France",
			OS:      "Android",
		}
		campaigns, err := deliveryStore.GetCampaigns(context.Background(), &query)
		fmt.Println(campaigns)
		assert.NoError(t, err)
		assert.Empty(t, campaigns)
	})

	t.Run("Should not return any campaigns if the AppID is in the exclusion list", func(t *testing.T) {

		query := model.Delivery{
			AppID:   "com.csg.game",
			Country: "US",
			OS:      "iOS",
		}
		campaigns, err := deliveryStore.GetCampaigns(context.Background(), &query)
		fmt.Println(campaigns)
		assert.NoError(t, err)
		assert.Empty(t, campaigns)
	})

	t.Run("Should not return any campaigns if the AppID is not in the inclusion list", func(t *testing.T) {

		query := model.Delivery{
			AppID:   "com.jetrun2.game",
			Country: "Germany",
			OS:      "iOS",
		}
		campaigns, err := deliveryStore.GetCampaigns(context.Background(), &query)
		fmt.Println(campaigns)
		assert.NoError(t, err)
		assert.Empty(t, campaigns)
	})
}
