package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Ruthvik10/targeting-engine/api"
	"github.com/Ruthvik10/targeting-engine/cache"
	"github.com/Ruthvik10/targeting-engine/config"
	"github.com/Ruthvik10/targeting-engine/store"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg, err := config.Load(".")
	if err != nil {
		log.Fatalf("Error loading the config: %v", err)
	}

	ctx, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()

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

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	campaignColl := client.Database("campaigns_db").Collection("campaigns")

	deliveryCache := cache.NewDelivery(redisClient)
	deliveryStore := store.NewDeliveryStore(campaignColl)
	deliveryHandler := api.NewDeliveryHandler(deliveryStore, deliveryCache)

	mux := http.NewServeMux()
	deliveryHandler.RegisterRoutes(mux)
	if err := http.ListenAndServe(cfg.ServerAddr, mux); err != nil {
		log.Fatalf("Unable to start the http server: %v", err)
	}
}
