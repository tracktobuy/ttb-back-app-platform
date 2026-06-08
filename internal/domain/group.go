package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Group struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UUID           string             `bson:"uuid" json:"uuid"`
	Version        int                `bson:"version" json:"version"`
	Name           string             `bson:"name" json:"name"`
	Budget         float32            `bson:"budget" json:"budget"`
	BudgetCurrency string             `bson:"budgetCurrency" json:"budgetCurrency"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
	CreatedBy      primitive.ObjectID `bson:"createdBy" json:"createdBy"`
}
