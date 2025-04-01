// src/hacker_rank/config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HackerRankAPIKey string
	Port             string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

	apiKey := os.Getenv("HACKERRANK_API_KEY")
	if apiKey == "" {
		return nil, log.Fatal("HACKERRANK_API_KEY environment variable not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		HackerRankAPIKey: apiKey,
		Port:             port,
	}, nil
}
