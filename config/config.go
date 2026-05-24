package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI        string
	MongoDB         string
	MongoServerPort string
}

func LoadConfig() *Config {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment...")
	}

	return &Config{
		MongoURI:        os.Getenv("MONGO_URI"),
		MongoDB:         os.Getenv("MONGO_DB"),
		MongoServerPort: os.Getenv("MONGO_SERVER_PORT"),
	}
}
