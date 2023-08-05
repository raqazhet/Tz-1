package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	AppPort    string `env:"APP_PORT" envDefault:"8000"`
	DBHost     string `env:"DB_HOST" envDefault:"localhost"`
	DBPort     string `env:"DB_PORT" envDefault:"27016"`
	DBUser     string `env:"DB_USER" envDefault:"region"`
	DBName     string `env:"DB_NAME" envDefault:"todo"`
	DBPassword string `env:"DB_PASSWORD"`
	TZ         string `env:"TZ" envDefault:"Asia/Almaty"`
}

func NewConfig() (*Config, error) {
	cfg := new(Config)
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("cannot parse env: %v", err)
	}
	return cfg, nil
}

func PrePareEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
}
