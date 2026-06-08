package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UUID      string             `bson:"uuid" json:"uuid"`
	Version   int                `bson:"version" json:"version"`
	Title     string             `bson:"title" json:"title"`
	Images    []string           `bson:"images" json:"images"`
	Labels    []string           `bson:"labels" json:"labels"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	CreatedBy primitive.ObjectID `bson:"createdBy" json:"createdBy"`
}
