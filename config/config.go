package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI      string
	MongoDB       string
	ApiServerPort string
}

func LoadConfig() *Config {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment...")
	}

	mongo_user := os.Getenv("MONGO_DB_USER")
	mongo_pass := os.Getenv("MONGO_DB_PASSWORD")
	mongo_host := os.Getenv("MONGO_HOST")
	mongo_db_name := os.Getenv("MONGO_DB_NAME")
	mongo_api_server_port := os.Getenv("API_SERVER_PORT")

	mongo_srv := "mongodb+srv://" + mongo_user + ":" + mongo_pass + "@" + mongo_host

	if mongo_user == "" || mongo_pass == "" || mongo_host == "" {
		log.Fatal("Missing required environment variables for MongoDB connection")
	}

	return &Config{
		MongoURI:      mongo_srv,
		MongoDB:       mongo_db_name,
		ApiServerPort: mongo_api_server_port,
	}
}
