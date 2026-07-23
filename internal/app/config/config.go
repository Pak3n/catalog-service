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

	if err := godotenv.Load(); err != nil {
		log.Printf("WARNING: .env file not found: %v", err)
	}

	if err := envconfig.Process("APP", &Root); err != nil {
		log.Fatal("failed to load config: ", err)
	}

}
