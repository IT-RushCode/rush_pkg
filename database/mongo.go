package database

import (
	"context"

	"github.com/IT-RushCode/rush_pkg/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoDBConnect(cfg *config.MongoDBConfig) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(cfg.URI)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func MongoDBClose(client *mongo.Client) error {
	err := client.Disconnect(context.Background())
	if err != nil {
		return err
	}

	return nil
}
