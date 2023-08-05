package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"tzregion/config"
	"tzregion/repository/mongodb"

	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// load env
	var once sync.Once
	once.Do(config.PrePareEnv)
	// get config
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}
	// init logger
	l, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	defer func(*zap.Logger) {
		err := l.Sync()
		if err != nil {
			log.Fatal(err)
		}
	}(l)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// repo layer
	client, err := mongodb.Dial(constructorUri(cfg))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	return nil
}

func constructorUri(cfg *config.Config) string {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	return uri
}
