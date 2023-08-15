package mongodb

import (
	"context"
	"fmt"
	"log"
	"tzregion/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Dial(ctx context.Context, cfg *config.Config, uri string) (*mongo.Client, error) {
	fmt.Println(uri)
	clientOptions := options.Client().ApplyURI(uri)
	clientOptions.SetAuth(options.Credential{
		Username: cfg.DBUser,
		Password: cfg.DBPassword,
	})
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Printf("mongo connect %v", err)
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Printf("ping doesn't work %v", err)
		return nil, err
	}
	return client, nil
}
