package response

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"-"`
	UUID     string             `json:"uuid"`
	Username string             `json:"username"`
	Name     string             `json:"name"`
	Groups   []Group            `json:"groups"`
}
