package store

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Ruthvik10/targeting-engine/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var deliveryStore *Delivery

func TestMain(m *testing.M) {

	cfg, err := config.Load("../")
	if err != nil {
		log.Fatalf("Error loading the config: %v", err)
	}
	ctx, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.TestDBSource))
	if err != nil {
		log.Fatalf("Error connecting to the database instance: %v", err)
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Could not ping MongoDB: ", err)
	}

	campaignColl := client.Database("campaigns_db").Collection("campaigns")
	deliveryStore = NewDelivery(campaignColl)
	os.Exit(m.Run())
}
