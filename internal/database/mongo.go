package database

import (
	"context"
	"log"
	"time"

	"github.com/SavanRajyaguru/ecommerce-go-notification-service/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Client *mongo.Client
var Collection *mongo.Collection

func ConnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := config.AppConfig.MongoURI
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to create Mongo client: %v", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB")
	Client = client
	// Default to "ecommerce_notifications" if not provided (fallback)
	dbName := config.AppConfig.MongoDBName
	if dbName == "" {
		dbName = "ecommerce_notifications"
	}
	Collection = client.Database(dbName).Collection("notifications")
}

func DisconnectDB() {
	if Client != nil {
		if err := Client.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting MongoDB: %v", err)
		}
	}
}
