package repository

import (
	"context"
	"time"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/response"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GroupRepository interface {
	Create(ctx context.Context, item domain.Group) (*domain.Group, error)
	Get(ctx context.Context, uuid string) (*domain.Group, error)
	Update(ctx context.Context, item *domain.Group) (*domain.Group, error)
	Delete(ctx context.Context, id string) error

	GetLabels(ctx context.Context, groupUUID string) (response.GroupLabel, error)
}

type mongoGroupRepo struct {
	collection   *mongo.Collection
	giCollection *mongo.Collection
}

func NewGroupRepo(db *mongo.Database) GroupRepository {
	return &mongoGroupRepo{
		collection:   db.Collection(GROUPS_COLLECTION_NAME),
		giCollection: db.Collection(GROUP_ITEM_COLLECTION_NAME),
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

func (g *mongoGroupRepo) Get(ctx context.Context, uuid string) (*domain.Group, error) {

	var group *domain.Group

	err := g.collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&group)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (g *mongoGroupRepo) Update(ctx context.Context, item *domain.Group) (*domain.Group, error) {

	version := item.Version

	item.UpdatedAt = time.Now().UTC()
	item.Version += 1

	update := bson.M{
		"$set": bson.M{
			"version":        item.Version,
			"name":           item.Name,
			"budget":         item.Budget,
			"budgetCurrency": item.BudgetCurrency,
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

func (g *mongoGroupRepo) GetLabels(ctx context.Context, groupUUID string) (response.GroupLabel, error) {

	group, err := g.Get(ctx, groupUUID)
	if err != nil {
		return response.GroupLabel{}, err
	}

	cursor, err := g.giCollection.Aggregate(ctx, g.labelsMongoPipeline(group.ID))
	if err != nil {
		return response.GroupLabel{}, err
	}

	defer cursor.Close(ctx)

	var result []response.GroupLabel
	if err = cursor.All(ctx, &result); err != nil {
		return response.GroupLabel{}, err
	}

	if len(result) > 0 {
		return result[0], nil
	}

	return response.GroupLabel{}, err
}

func (g *mongoGroupRepo) labelsMongoPipeline(groupID primitive.ObjectID) mongo.Pipeline {

	return mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"groupId": groupID}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "items",
			"localField":   "itemId",
			"foreignField": "_id",
			"as":           "itemDetails",
		}}},

		{{Key: "$lookup", Value: bson.M{
			"from":         "groups",
			"localField":   "groupId",
			"foreignField": "_id",
			"as":           "groupDetails",
		}}},

		{{Key: "$unwind", Value: "$itemDetails"}},
		{{Key: "$unwind", Value: "$itemDetails.labels"}},
		{{Key: "$unwind", Value: "$groupDetails"}},

		{{Key: "$group", Value: bson.M{
			"_id":       nil,
			"items":     bson.M{"$addToSet": "$itemDetails"},
			"labels":    bson.M{"$addToSet": "$itemDetails.labels"},
			"groupName": bson.M{"$first": "$groupDetails.uuid"},
			"groupUUID": bson.M{"$first": "$groupDetails.uuid"},
		}}},

		{{Key: "$project", Value: bson.M{
			"_id":       0,
			"items":     1,
			"labels":    1,
			"groupName": 1,
			"groupUUID": 1,
		}}},
	}
}
