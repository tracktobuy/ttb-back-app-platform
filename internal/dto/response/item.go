package response

import "go.mongodb.org/mongo-driver/bson/primitive"

type Item struct {
	ID        primitive.ObjectID `json:"-"`
	UUID      string             `json:"uuid"`
	Title     string             `json:"title"`
	Images    []string           `json:"images,omitempty"`
	CreatedAt string             `json:"createdAt"`
	UpdatedAt string             `json:"updatedAt,omitempty"`
	CreatedBy User               `json:"createdBy"`
	Labels    []string           `json:"labels,omitempty"`
	Stores    []string           `json:"stores,omitempty"`
	Groups    []string           `json:"groups,omitempty"`
}
