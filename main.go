package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Ruthvik10/targeting-engine/api"
	"github.com/Ruthvik10/targeting-engine/config"
	"github.com/Ruthvik10/targeting-engine/store"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg, err := config.Load(".")
	if err != nil {
		log.Fatalf("Error loading the config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DBURI))
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

	log.Println("Successfully connected to the database instance")

	campaignColl := client.Database("campaigns_db").Collection("campaigns")

	deliveryStore := store.NewDeliveryStore(campaignColl)
	deliveryHandler := api.NewDeliveryHandler(deliveryStore)

	mux := http.NewServeMux()
	deliveryHandler.RegisterRoutes(mux)
	if err := http.ListenAndServe(cfg.ServerAddr, mux); err != nil {
		log.Fatalf("Unable to start the http server: %v", err)
	}
}
