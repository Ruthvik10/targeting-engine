package store

import (
	"context"
	"log"

	"github.com/Ruthvik10/targeting-engine/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Delivery struct {
	collection *mongo.Collection
}

func NewDelivery(coll *mongo.Collection) *Delivery {
	return &Delivery{collection: coll}
}

func (store *Delivery) GetCampaigns(ctx context.Context, delivery *model.Delivery) ([]*model.Campaign, error) {
	filter := bson.M{
		"status": "ACTIVE",
		"$and": []bson.M{
			{
				"$or": []bson.M{
					{"targeting.includeOS": delivery.OS},
					{"targeting.includeOS": bson.M{"$size": 0}},
				},
			},
			{
				"$or": []bson.M{
					{"targeting.includeCountry": delivery.Country},
					{"targeting.includeCountry": bson.M{"$size": 0}},
				},
			},
			{
				"$or": []bson.M{
					{"targeting.includeApp": delivery.AppID},
					{"targeting.includeApp": bson.M{"$size": 0}},
				},
			},
		},
		"targeting.excludeOS":      bson.M{"$nin": []string{delivery.OS}},
		"targeting.excludeCountry": bson.M{"$nin": []string{delivery.Country}},
		"targeting.excludeApp":     bson.M{"$nin": []string{delivery.AppID}},
	}

	findOptions := &options.FindOptions{
		Collation: &options.Collation{Locale: "en", Strength: 2},
	}

	cursor, err := store.collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Printf("Error fetching the campaigns from the database: %v\n", err)
		return nil, err
	}
	var campaigns []*model.Campaign

	// Iterate over the cursor and decode each document into a Campaign struct
	for cursor.Next(ctx) {
		var campaign model.Campaign
		if err := cursor.Decode(&campaign); err != nil {
			log.Printf("Error decoding the campaign cursor values: %v\n", err)
			return nil, err
		}
		campaigns = append(campaigns, &campaign)
	}

	// Check if any error occurred during iteration
	if err := cursor.Err(); err != nil {
		log.Printf("Error fetching the campaigns from the database: %v\n", err)
		return nil, err
	}
	return campaigns, nil
}

func (store *Delivery) WatchCampaign(ctx context.Context, out chan<- bson.M) {

	// Only listen to update operations on the campaign collection and fetch the document only if the status has changed to INACTIVE
	matchPipeline := bson.D{
		{
			Key: "$match", Value: bson.D{
				{Key: "operationType", Value: "update"},
				{Key: "fullDocument.status", Value: "INACTIVE"},
			},
		},
	}
	csOptions := options.ChangeStream().SetFullDocument(options.UpdateLookup)
	cs, err := store.collection.Watch(ctx, mongo.Pipeline{matchPipeline}, csOptions)
	if err != nil {
		log.Printf("Error watching the campaign collection for changes: %v\n", err)
	}
	defer cs.Close(ctx)
	for cs.Next(ctx) {
		var data bson.M
		if err := cs.Decode(&data); err != nil {
			panic(err)
		}
		out <- data
	}
}
