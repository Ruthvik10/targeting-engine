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

	log.Println("Successfully connected to the database instance.")

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisURI,
	})

	log.Println("Connected to the cache.")

	campaignColl := client.Database("campaigns_db").Collection("campaigns")

	cacheExpiry := time.Duration(cfg.RedisCacheExpr) * time.Second
	deliveryCache := cache.NewDelivery(redisClient)
	deliveryStore := store.NewDelivery(campaignColl)
	deliveryHandler := api.NewDeliveryHandler(deliveryStore, deliveryCache, cacheExpiry)

	mux := http.NewServeMux()
	deliveryHandler.RegisterRoutes(mux)

	log.Println("Starting the http server on: " + cfg.ServerAddr)
	if err := http.ListenAndServe(cfg.ServerAddr, mux); err != nil {
		log.Fatalf("Unable to start the http server: %v", err)
	}
}
