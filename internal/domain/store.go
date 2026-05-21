package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Store struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	UUID         string             `bson:"uuid" json:"uuid"`
	Version      int                `bson:"version" json:"-"`
	Price        float32            `bson:"price" json:"price"`
	ShippingCost float32            `bson:"shippingCost" json:"shippingCost"`
	Currency     string             `bson:"currency" json:"currency"`
	Domain       string             `bson:"domain" json:"domain"`
	Name         string             `bson:"name" json:"name"`
	BestOption   bool               `bson:"bestOption" json:"bestOption"`
	URL          string             `bson:"url" json:"url"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
	CreatedBy    primitive.ObjectID `bson:"createdBy" json:"-"`
	Item         primitive.ObjectID `bson:"item" json:"-"`
}
