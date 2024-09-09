package store

import (
	"context"

	"github.com/Ruthvik10/targeting-engine/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Delivery struct {
	collection *mongo.Collection
}

func NewDeliveryStore(coll *mongo.Collection) *Delivery {
	return &Delivery{collection: coll}
}

func (store *Delivery) GetCampaigns(ctx context.Context, in *model.Delivery) ([]*model.Campaign, error) {
	filter := bson.M{
		"status": "ACTIVE",
		"$and": []bson.M{
			{
				"$or": []bson.M{
					{"targeting.includeOS": in.OS},
					{"targeting.includeOS": bson.M{"$size": 0}},
				},
			},
			{
				"$or": []bson.M{
					{"targeting.includeCountry": in.Country},
					{"targeting.includeCountry": bson.M{"$size": 0}},
				},
			},
			{
				"$or": []bson.M{
					{"targeting.includeApp": in.AppID},
					{"targeting.includeApp": bson.M{"$size": 0}},
				},
			},
		},
		"targeting.excludeOS":      bson.M{"$nin": []string{in.OS}},
		"targeting.excludeCountry": bson.M{"$nin": []string{in.Country}},
		"targeting.excludeApp":     bson.M{"$nin": []string{in.AppID}},
	}

	findOptions := &options.FindOptions{
		Collation: &options.Collation{Locale: "en", Strength: 2},
	}

	cursor, err := store.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	var campaigns []*model.Campaign

	// Iterate over the cursor and decode each document into a Campaign struct
	for cursor.Next(ctx) {
		var campaign model.Campaign
		if err := cursor.Decode(&campaign); err != nil {
			return nil, err
		}
		campaigns = append(campaigns, &campaign)
	}

	// Check if any error occurred during iteration
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return campaigns, nil
}
