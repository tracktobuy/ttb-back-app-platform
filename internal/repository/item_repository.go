package repository

import (
	"context"
	"time"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"github.com/tracktobuy/ttb-back-app-platform/internal/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ItemRepository interface {
	Create(ctx context.Context, item domain.Item) (*domain.Item, error)
	Get(ctx context.Context, itemUUID string) (*domain.Item, error)
	Delete(ctx context.Context, id string) error
	GetAllByGroupID(ctx context.Context, groupID primitive.ObjectID) ([]domain.Item, error)
	GetAllByUserId(ctx context.Context, userId primitive.ObjectID) ([]domain.Item, error)
}

type itemRepo struct {
	log          logger.Logger
	collection   *mongo.Collection
	ugCollection *mongo.Collection
}

type aggregationResult struct {
	ItemsDetails []domain.Item `bson:"items_details"`
}

func NewItemRepo(db *mongo.Database) ItemRepository {
	log := logger.NewLogger()
	log.SetRepositoryName("ItemRepository")
	return &itemRepo{
		collection:   db.Collection(ITEMS_COLLECTION_NAME),
		ugCollection: db.Collection(USER_GROUP_COLLECTION_NAME),
		log:          log,
	}
}

func (r *itemRepo) Create(ctx context.Context, item domain.Item) (*domain.Item, error) {
	item.UUID = helper.GenerateUUIDV7()
	item.Version = 1
	item.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *itemRepo) Get(ctx context.Context, itemUUID string) (*domain.Item, error) {
	var item domain.Item
	err := r.collection.FindOne(ctx, bson.M{"uuid": itemUUID}).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *itemRepo) Delete(ctx context.Context, id string) error {

	_, err := r.collection.DeleteOne(ctx, bson.M{"uuid": id})
	if err != nil {
		return err
	}

	return nil
}

func (r *itemRepo) GetAllByGroupID(ctx context.Context, groupID primitive.ObjectID) ([]domain.Item, error) {
	var items []domain.Item

	filterList := []primitive.ObjectID{groupID}

	filter := bson.M{
		"$in": filterList,
	}

	cursor, err := r.collection.Find(ctx, bson.M{"groups": filter})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &items); err != nil {
		return []domain.Item{}, err
	}

	return items, nil
}

func (r *itemRepo) GetAllByUserId(ctx context.Context, userId primitive.ObjectID) ([]domain.Item, error) {

	r.log.SetMethodName("GetAllByUserId")

	cursor, err := r.ugCollection.Aggregate(ctx, r.getPipeline(userId))
	if err != nil {
		return []domain.Item{}, err
	}

	defer cursor.Close(ctx)

	var results []aggregationResult

	if err = cursor.All(ctx, &results); err != nil {
		r.log.Error("failed to decode results", "error", err.Error())
		return []domain.Item{}, err
	}

	if len(results) == 0 {
		r.log.Error("no results found", "userId", userId, "total items found", len(results))
		return []domain.Item{}, nil
	}

	return results[0].ItemsDetails, nil

}

func (r *itemRepo) getPipeline(userId primitive.ObjectID) mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"userId": userId}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "userId",
			"foreignField": "_id",
			"as":           "users_details",
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "groups",
			"localField":   "groupId",
			"foreignField": "_id",
			"as":           "group_details",
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "group_item",
			"localField":   "groupId",
			"foreignField": "groupId",
			"as":           "group_item",
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "items",
			"localField":   "group_item.itemId",
			"foreignField": "_id",
			"as":           "items_details",
		}}},
		{{Key: "$unwind", Value: "$items_details"}},
		{{Key: "$group", Value: bson.M{
			"_id":           nil,
			"items_details": bson.M{"$addToSet": "$items_details"},
		}}},
		{{Key: "$project", Value: bson.M{
			"_id":           0,
			"items_details": 1,
		}}},
	}
}
