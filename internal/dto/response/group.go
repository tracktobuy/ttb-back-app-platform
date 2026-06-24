package response

import "go.mongodb.org/mongo-driver/bson/primitive"

type Group struct {
	ID             primitive.ObjectID `json:"-"`
	UUID           string             `json:"uuid"`
	Name           string             `json:"name"`
	Budget         float32            `json:"budget"`
	BudgetCurrency string             `json:"budgetCurrency"`
	CreatedAt      string             `json:"createdAt"`
	UpdatedAt      string             `json:"updatedAt,omitempty"`
	CreatedBy      User               `json:"createdBy,omitempty"`
}
