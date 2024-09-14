package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Ruthvik10/targeting-engine/jsonutil"
	"github.com/Ruthvik10/targeting-engine/model"
)

type store interface {
	GetCampaigns(ctx context.Context, in *model.Delivery) ([]*model.Campaign, error)
}

type cache interface {
	GetCampaigns(ctx context.Context, key string) (string, error)
	SetCampaign(ctx context.Context, key, value string, exp time.Duration) error
}

type DeliveryHandler struct {
	store       store
	cache       cache
	cacheExpiry time.Duration
}

func NewDeliveryHandler(store store, cache cache, cacheExpiry time.Duration) *DeliveryHandler {
	return &DeliveryHandler{store: store, cache: cache, cacheExpiry: cacheExpiry}
}

func (h *DeliveryHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /delivery", h.DeliverCampaign)
}

type DeliveryResponse struct {
	CID   string `json:"cid"`
	Image string `json:"img"`
	CTA   string `json:"cta"`
}

func (h *DeliveryHandler) DeliverCampaign(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var app, country, os string

	if len(queryParams["app"]) == 0 {
		log.Println("Missing app param in the request")
		jsonutil.WriteError(w, http.StatusBadRequest, "missing app param")
		return
	}

	if len(queryParams["country"]) == 0 {
		log.Println("Missing country param in the request")
		jsonutil.WriteError(w, http.StatusBadRequest, "missing country param")
		return
	}

	if len(queryParams["os"]) == 0 {
		log.Println("Missing os param in the request")
		jsonutil.WriteError(w, http.StatusBadRequest, "missing os param")
		return
	}

	app = queryParams["app"][0]
	country = queryParams["country"][0]
	os = queryParams["os"][0]

	cacheKey := fmt.Sprintf("%s:%s:%s", app, country, os)

	cachedValue, _ := h.cache.GetCampaigns(context.Background(), cacheKey)

	campaigns := make([]*model.Campaign, 0)
	if cachedValue == "" {
		log.Printf("Cache miss for the key: (%s), quering the database\n", cacheKey)
		var err error
		campaigns, err = h.store.GetCampaigns(context.Background(), &model.Delivery{AppID: app, Country: country, OS: os})
		if err != nil {
			jsonutil.WriteError(w, http.StatusInternalServerError, "Something went wrong while fetching the campaigns!")
			return
		}

		if len(campaigns) == 0 {
			jsonutil.WriteError(w, http.StatusNotFound, "No campaign exists for the specified parameters!")
			return
		}

		// Store the fetched data in cache.
		// If there is any error, do not throw the error to the client since the requested is already available from the database.
		bytes, _ := json.Marshal(campaigns)
		err = h.cache.SetCampaign(context.Background(), cacheKey, string(bytes), h.cacheExpiry)
		if err == nil {
			log.Printf("Added value for the key: (%s) to the cache\n", cacheKey)
		}
	} else {
		log.Printf("Cache hit for the key: (%s)\n", cacheKey)
		err := json.Unmarshal([]byte(cachedValue), &campaigns)
		if err != nil {
			jsonutil.WriteError(w, http.StatusInternalServerError, "Something went wrong!")
			return
		}
	}

	out := make([]*DeliveryResponse, 0)

	for _, c := range campaigns {
		out = append(out, &DeliveryResponse{
			CID:   c.ID.Hex(),
			Image: c.Image,
			CTA:   c.CTA,
		})
	}

	jsonutil.WriteJSON(w, http.StatusOK, out)
}
