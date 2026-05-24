package repository

import (
	"context"
	"time"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoUserRepo struct {
	collection *mongo.Collection
}

func NewUserRepo(db *mongo.Database) CrudRepository[domain.User] {
	return &mongoUserRepo{
		collection: db.Collection("users"),
	}
}

func (u *mongoUserRepo) Create(ctx context.Context, item domain.User) (*domain.User, error) {

	item.ID = primitive.NewObjectID()
	item.CreatedAt = time.Now()

	_, err := u.collection.InsertOne(ctx, item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (u *mongoUserRepo) Get(ctx context.Context, id string) (*domain.User, error) {

	var user *domain.User

	err := u.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *mongoUserRepo) GetAll(ctx context.Context) ([]domain.User, error) {

	cursor, err := u.collection.Find(ctx, bson.M{})

	if err != nil {
		return []domain.User{}, err
	}

	var users []domain.User

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &users); err != nil {
		return []domain.User{}, err
	}

	if len(users) == 0 {
		return []domain.User{}, nil
	}

	return users, nil
}

func (u *mongoUserRepo) Update(ctx context.Context, item domain.User) (*domain.User, error) {

	version := item.Version

	item.UpdatedAt = time.Now()
	item.Version += 1

	update := bson.M{
		"$set": bson.M{
			"version":   item.Version,
			"name":      item.Name,
			"username":  item.Username,
			"groups":    item.Groups,
			"updatedAt": item.UpdatedAt,
		},
	}

	_, err := u.collection.UpdateOne(ctx, bson.M{"_id": item.ID, "version": version}, update)
	if err != nil {
		return nil, err
	}

	return u.Get(ctx, item.ID.Hex())
}

func (u *mongoUserRepo) Delete(ctx context.Context, id string) error {
	if _, err := u.collection.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		return err
	}

	return nil
}
