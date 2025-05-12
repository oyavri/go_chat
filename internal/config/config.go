package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	Hostname        string
	DbConnectionUrl string
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetConfig() *Config {
	DbPort := os.Getenv("DB_PORT")
	DbUsername := os.Getenv("DB_USERNAME")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbHostname := os.Getenv("DB_HOSTNAME")
	DbNameUser := os.Getenv("DB_NAME")

	baseConnUrl := fmt.Sprintf("postgres://%v:%v@%v:%v/", DbUsername, DbPassword, DbHostname, DbPort)
	connectionUrl := baseConnUrl + DbNameUser

	return &Config{
		Port:            os.Getenv("PORT"),
		Hostname:        os.Getenv("HOSTNAME"),
		DbConnectionUrl: connectionUrl,
	}
}
