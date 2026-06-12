package response

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Store struct {
	ID           primitive.ObjectID `json:"-"`
	UUID         string             `json:"uuid"`
	Price        float32            `json:"price"`
	ShippingCost float32            `json:"shippingCost"`
	Currency     string             `json:"currency"`
	Domain       string             `json:"domain"`
	Name         string             `json:"name"`
	BestOption   bool               `json:"bestOption"`
	URL          string             `json:"url"`
	CreatedAt    string             `json:"createdAt"`
	UpdatedAt    string             `json:"updatedAt,omitempty"`
	CreatedBy    string             `json:"createdBy,omitempty"`
	ItemId       string             `json:"itemId,omitempty"`
}
