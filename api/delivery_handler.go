package api

import (
	"context"
	"net/http"

	"github.com/Ruthvik10/targeting-engine/jsonutil"
	"github.com/Ruthvik10/targeting-engine/model"
)

type CampaignStore interface {
	GetCampaigns(ctx context.Context, in *model.Delivery) ([]*model.Campaign, error)
}

type DeliveryHandler struct {
	store CampaignStore
}

func NewDeliveryHandler(store CampaignStore) *DeliveryHandler {
	return &DeliveryHandler{store: store}
}

func (h *DeliveryHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /delivery", h.DeliverCampaign)
}

func (h *DeliveryHandler) DeliverCampaign(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	app := queryParams["app"][0]
	country := queryParams["country"][0]
	os := queryParams["os"][0]

	// Handle bad request
	campaigns, err := h.store.GetCampaigns(context.Background(), &model.Delivery{AppID: app, Country: country, OS: os})
	if err != nil {
		jsonutil.WriteError(w, http.StatusInternalServerError, "Something went wrong while fetching the campaigns!")
		return
	}
	jsonutil.WriteJSON(w, http.StatusOK, campaigns)
}
