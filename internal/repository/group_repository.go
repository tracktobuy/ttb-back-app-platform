package repository

import (
	"context"
	"time"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GroupRepositoryInterface interface {
	CrudRepository[domain.Group]
	GetByUUID(ctx context.Context, uuid string) (*domain.Group, error)
}

type mongoGroupRepo struct {
	collection *mongo.Collection
}

func NewGroupRepo(db *mongo.Database) GroupRepositoryInterface {
	return &mongoGroupRepo{
		collection: db.Collection("groups"),
	}
}

func (g *mongoGroupRepo) Create(ctx context.Context, item domain.Group) (*domain.Group, error) {

	item.ID = primitive.NewObjectID()
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

	err := g.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&group)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (g *mongoGroupRepo) GetAll(ctx context.Context) ([]domain.Group, error) {

	cursor, err := g.collection.Find(ctx, bson.M{})

	if err != nil {
		return []domain.Group{}, err
	}

	var groups []domain.Group

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &groups); err != nil {
		return []domain.Group{}, err
	}

	if len(groups) == 0 {
		return []domain.Group{}, nil
	}

	return groups, nil
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

	return g.Get(ctx, item.ID.Hex())
}

func (g *mongoGroupRepo) Delete(ctx context.Context, id string) error {
	if _, err := g.collection.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		return err
	}

	return nil
}

func (g *mongoGroupRepo) GetByUUID(ctx context.Context, uuid string) (*domain.Group, error) {
	var group *domain.Group

	err := g.collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&group)
	if err != nil {
		return nil, err
	}

	return group, nil
}
