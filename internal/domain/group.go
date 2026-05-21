package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Group struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty" json:"-"`
	UUID           string               `bson:"uuid" json:"uuid"`
	Version        int                  `bson:"version" json:"-"`
	Name           string               `bson:"name" json:"name"`
	Budget         float32              `bson:"budget" json:"budget"`
	BudgetCurrency string               `bson:"budgetCurrency" json:"budgetCurrency"`
	Items          []primitive.ObjectID `bson:"items" json:"-"`
	Users          []primitive.ObjectID `bson:"users" json:"-"`
	CreatedAt      time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time            `bson:"updatedAt" json:"updatedAt"`
	CreatedBy      primitive.ObjectID   `bson:"createdBy" json:"-"`
}
