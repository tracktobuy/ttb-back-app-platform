package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"-"`
	UUID      string               `bson:"uuid" json:"uuid"`
	Version   int                  `bson:"version" json:"-"`
	Title     string               `bson:"title" json:"title"`
	Images    []string             `bson:"images" json:"images"`
	Labels    []string             `bson:"labels" json:"labels"`
	Stores    []primitive.ObjectID `bson:"stores" json:"-"`
	Groups    []primitive.ObjectID `bson:"groups" json:"-"`
	CreatedAt time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time            `bson:"updatedAt" json:"updatedAt"`
	CreatedBy primitive.ObjectID   `bson:"createdBy" json:"-"`
}
