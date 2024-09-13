package storage

import (
	"context"

	"github.com/IT-RushCode/rush_pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MONGO_DB_CONNECT(ctx context.Context, cfg *config.MongoDBConfig) *mongo.Client {
	clientOptions := options.Client().ApplyURI(cfg.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic("mongoDB is not connected!")
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		panic("mongoDB is not connected!")
	}

	return client
}

func MONGO_DB_CLOSE(ctx context.Context, client *mongo.Client) error {
	err := client.Disconnect(ctx)
	if err != nil {
		return err
	}

	return nil
}
