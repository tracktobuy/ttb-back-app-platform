package repository

import (
	"context"
	"time"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GroupRepository interface {
	Create(ctx context.Context, item domain.Group) (*domain.Group, error)
	Get(ctx context.Context, id string) (*domain.Group, error)
	Update(ctx context.Context, item domain.Group) (*domain.Group, error)
	Delete(ctx context.Context, id string) error
}

type mongoGroupRepo struct {
	collection *mongo.Collection
}

func NewGroupRepo(db *mongo.Database) GroupRepository {
	return &mongoGroupRepo{
		collection: db.Collection("groups"),
	}
}

func (g *mongoGroupRepo) Create(ctx context.Context, item domain.Group) (*domain.Group, error) {

	item.ID = primitive.NewObjectID()
	item.UUID = helper.GenerateUUIDV7()
	item.Version = 1
	item.CreatedAt = time.Now().UTC()

	_, err := g.collection.InsertOne(ctx, item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (g *mongoGroupRepo) Get(ctx context.Context, id string) (*domain.Group, error) {

	var group *domain.Group

	err := g.collection.FindOne(ctx, bson.M{"uuid": id}).Decode(&group)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (g *mongoGroupRepo) Update(ctx context.Context, item domain.Group) (*domain.Group, error) {

	version := item.Version

	item.UpdatedAt = time.Now().UTC()
	item.Version += 1

	update := bson.M{
		"$set": bson.M{
			"version":        item.Version,
			"name":           item.Name,
			"budget":         item.Budget,
			"budgetCurrency": item.BudgetCurrency,
			"items":          item.Items,
			"users":          item.Users,
			"updatedAt":      item.UpdatedAt,
		},
	}

	_, err := g.collection.UpdateOne(ctx, bson.M{"_id": item.ID, "version": version}, update)
	if err != nil {
		return nil, err
	}

	return g.Get(ctx, item.UUID)
}

func (g *mongoGroupRepo) Delete(ctx context.Context, id string) error {
	if _, err := g.collection.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		return err
	}

	return nil
}
