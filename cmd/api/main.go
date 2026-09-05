package main

import (
	"context"
	"time"

	"github.com/tracktobuy/ttb-back-app-platform/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type application struct {
	handlers *Handler
}

func main() {

	cfg := config.LoadConfig()

	db := LoadDatabase(cfg)
	defer db.Client().Disconnect(context.Background())

	app := &application{
		handlers: CreateHandlers(db),
	}

	app.run(cfg)
}

func LoadDatabase(cfg *config.Config) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := config.MongoConnect(ctx, cfg)

	return client.Database(cfg.MongoDB)
}
