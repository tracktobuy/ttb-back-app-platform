package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Origin used when CORS_ALLOWED_ORIGINS is not set, matching the local
// frontend development server.
//const defaultCorsAllowedOrigin = "http://localhost:3000"

type Config struct {
	MongoURI           string
	MongoDB            string
	ApiServerPort      string
	CorsAllowedOrigins []string
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
	cors_allowed_origins := os.Getenv("CORS_ALLOWED_ORIGINS")

	mongo_srv := "mongodb+srv://" + mongo_user + ":" + mongo_pass + "@" + mongo_host

	if mongo_user == "" || mongo_pass == "" || mongo_host == "" {
		log.Fatal("Missing required environment variables for MongoDB connection")
	}

	return &Config{
		MongoURI:           mongo_srv,
		MongoDB:            mongo_db_name,
		ApiServerPort:      mongo_api_server_port,
		CorsAllowedOrigins: parseCorsAllowedOrigins(cors_allowed_origins),
	}
}

// parseCorsAllowedOrigins splits a comma separated list of origins, falling
// back to the local frontend when none are configured.
func parseCorsAllowedOrigins(raw string) []string {

	origins := []string{}

	for _, origin := range strings.Split(raw, ",") {
		if trimmed := strings.TrimSpace(origin); trimmed != "" {
			origins = append(origins, trimmed)
		}
	}

	// if len(origins) == 0 {
	// 	return []string{defaultCorsAllowedOrigin}
	// }

	return origins
}
