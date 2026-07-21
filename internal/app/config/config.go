package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/Pak3n/catalog-service/internal/app/config/section"
)

type Config struct {
	Repository section.Repository `env:"REPOSITORY"`
	Processor  section.Processor  `env:"PROCESSOR"`
	Monitor    section.Monitor    `env:"MONITOR"`
}

var Root Config

func Load() {
	_ = godotenv.Load()

	if err := envconfig.Process("APP", &Root); err != nil {
		log.Fatal("failed to load config: ", err)
	}
	if Root.Repository.Postgres.Address == "" {
		log.Fatal("APP_REPOSITORY_POSTGRES_ADDRESS is required")
	}
	if Root.Repository.Postgres.Username == "" {
		log.Fatal("APP_REPOSITORY_POSTGRES_USERNAME is required")
	}
	if Root.Repository.Postgres.Password == "" {
		log.Fatal("APP_REPOSITORY_POSTGRES_PASSWORD is required")
	}
}
