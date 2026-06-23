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
	GetItems(ctx context.Context, groupUUID string) (response.GroupItem, error)
}

type mongoGroupRepo struct {
	collection   *mongo.Collection
	giCollection *mongo.Collection
}

type tmpResult struct {
	GroupUUID string    `bson:"groupUUID"`
	GroupName string    `bson:"groupName"`
	Items     []tmpItem `bson:"items"`
}

type tmpItem struct {
	ID          primitive.ObjectID `bson:"_id"`
	UUID        string             `bson:"uuid"`
	Title       string             `bson:"title"`
	Images      []string           `bson:"images"`
	Labels      []string           `bson:"labels"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
	LowestPrice float32            `bson:"lowestPrice"`
	CreatedBy   response.User      `bson:"createdBy"`
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

func (g *mongoGroupRepo) GetItems(ctx context.Context, groupUUID string) (response.GroupItem, error) {

	group, err := g.Get(ctx, groupUUID)
	if err != nil {
		return response.GroupItem{}, err
	}

	cursor, err := g.giCollection.Aggregate(ctx, g.itemsMongoPipeline(group.ID))
	if err != nil {
		return response.GroupItem{}, err
	}

	defer cursor.Close(ctx)

	var tmpResult []tmpResult
	if err = cursor.All(ctx, &tmpResult); err != nil {
		return response.GroupItem{}, err
	}

	if len(tmpResult) == 0 {
		return response.GroupItem{}, err
	}

	items := tmpResult[0].Items
	var respItems []response.Item

	for _, i := range items {
		tmp := response.Item{
			ID:          i.ID,
			UUID:        i.UUID,
			Title:       i.Title,
			Images:      i.Images,
			CreatedAt:   helper.DateTime(i.CreatedAt),
			UpdatedAt:   helper.DateTime(i.UpdatedAt),
			Labels:      i.Labels,
			LowestPrice: i.LowestPrice,
			CreatedBy: response.User{
				UUID:     i.CreatedBy.UUID,
				Name:     i.CreatedBy.Name,
				Username: i.CreatedBy.Username,
			},
		}
		respItems = append(respItems, tmp)
	}

	return response.GroupItem{
		GroupUUID: tmpResult[0].GroupUUID,
		GroupName: tmpResult[0].GroupName,
		Items:     respItems,
	}, nil

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

func (g *mongoGroupRepo) itemsMongoPipeline(groupID primitive.ObjectID) mongo.Pipeline {

	return mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"groupId": groupID}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "items",
			"localField":   "itemId",
			"foreignField": "_id",
			"as":           "itemDetails",
		}}},

		{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "itemDetails.createdBy",
			"foreignField": "_id",
			"as":           "user",
		}}},

		{{Key: "$lookup", Value: bson.M{
			"from":         "groups",
			"localField":   "groupId",
			"foreignField": "_id",
			"as":           "groupDetails",
		}}},

		{{Key: "$unwind", Value: "$itemDetails"}},
		{{Key: "$unwind", Value: "$groupDetails"}},
		{{Key: "$unwind", Value: "$user"}},

		{{Key: "$lookup", Value: bson.M{
			"from":         "stores",
			"localField":   "itemId",
			"foreignField": "itemId",
			"as":           "storeOptions",
		}}},

		{{Key: "$group", Value: bson.M{
			"_id":       "$groupId",
			"groupName": bson.M{"$first": "$groupDetails.name"},
			"groupUUID": bson.M{"$first": "$groupDetails.uuid"},
			"items": bson.M{
				"$push": bson.M{
					"uuid":        "$itemDetails.uuid",
					"title":       "$itemDetails.title",
					"images":      "$itemDetails.images",
					"labels":      "$itemDetails.labels",
					"createdAt":   "$itemDetails.createdAt",
					"lowestPrice": bson.M{"$min": "$storeOptions.price"},
					"createdBy": bson.M{
						"uuid":     "$user.uuid",
						"name":     "$user.name",
						"username": "$user.username",
					},
				},
			},
		}}},

		{{Key: "$project", Value: bson.M{
			"_id":       0,
			"groupName": 1,
			"groupUUID": 1,
			"items":     1,
		}}},
	}
}
