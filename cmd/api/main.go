package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/tracktobuy/ttb-back-app-platform/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type application struct {
	handlers *Handler
}

func main() {

	db := LoadDatabase()
	defer db.Client().Disconnect(context.Background())

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	app := &application{
		handlers: CreateHandlers(db).WithLogger(logger),
	}

	app.run()
}

func LoadDatabase() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := config.LoadConfig()

	client := config.MongoConnect(ctx, cfg)

	return client.Database(cfg.MongoDB)
}
