package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"tzregion/config"
	"tzregion/repository"
	"tzregion/repository/mongodb"
	"tzregion/service"
	"tzregion/transport/http"
	"tzregion/transport/http/handler"

	"go.uber.org/zap"
)

// @title  QR
// @version 1.0
// @description API server for todolist Application

// @host localhost:8000
// @BasePath /

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
	ctxx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongodb.Dial(ctxx, cfg, constructorUri(cfg))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctxx)
	db := client.Database(cfg.DBName)
	repotodo := repository.NewRepository(db)
	// service layer
	serviceTodo := service.NewService(repotodo, l)
	// handler layer
	handler := handler.NewHandler(serviceTodo, l)
	// http server instance
	httpServer := http.NewServer(cfg, handler)
	l.Info("Start app", zap.String("port", cfg.AppPort))
	httpServer.StartServer()

	// grace full shutdown
	quite := make(chan os.Signal, 1)
	signal.Notify(quite, syscall.SIGINT, syscall.SIGTERM)
	select {
	case s := <-quite:
		l.Info("signal accepted: ", zap.String("signal", s.String()))
	case err := <-httpServer.Notify:
		l.Info("server closing", zap.Error(err))
	}
	if err := httpServer.Shutdown(); err != nil {
		return fmt.Errorf("error while shutting down server: %s", err)
	}
	return nil
}

func constructorUri(cfg *config.Config) string {
	uri := fmt.Sprintf("mongodb://%s:%s", cfg.DBHost, cfg.DBPort)
	return uri
}
