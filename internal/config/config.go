package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	Port string
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetConfig() *config {
	return &config{
		Port: os.Getenv("PORT"),
	}
}
