package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	MongoURI    string
	MongoDBName string
}

func NewConfig() *Config {
	// Attempt to load .env (optional)
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, or error loading it; proceeding...")
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		Port:        port,
		MongoURI:    os.Getenv("MONGO_URI"),
		MongoDBName: os.Getenv("MONGO_DB"),
	}
}
