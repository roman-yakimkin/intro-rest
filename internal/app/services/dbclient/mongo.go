package dbclient

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"intro-rest/internal/app/services/configmanager"
	"log"
)

type MongoDBClient struct {
	config *configmanager.Config
	client *mongo.Client
}

func NewMongoDBClient(config *configmanager.Config) *MongoDBClient {
	return &MongoDBClient{
		config: config,
	}
}

func (c *MongoDBClient) Connect(ctx context.Context) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(c.config.MongoURI)
	var err error
	c.client, err = mongo.Connect(ctx, clientOptions)
	return c.client, err
}

func (c *MongoDBClient) Disconnect(ctx context.Context) {
	if err := c.client.Disconnect(ctx); err != nil {
		log.Fatal(err)
	}
}
