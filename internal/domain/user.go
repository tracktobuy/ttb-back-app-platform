package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UUID      string             `bson:"uuid" json:"uuid"`
	Version   int                `bson:"version" json:"version"`
	Name      string             `bson:"name" json:"name"`
	Username  string             `bson:"username" json:"username"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt,omitempty"`
}
