package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Dial(uri string) (*mongo.Database, error) {
	// Set up a context to control the connection's lifetime
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Set up options for the MongoDB client
	clientOptions := options.Client().ApplyURI(uri)
	// connect to mongo
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)
	// Check if the connection was successful
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot connect ot redis: %w", err)
	}
	db := client.Database("todo-list")
	return db, nil
}
