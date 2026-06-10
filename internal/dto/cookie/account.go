package cookie

import (
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Account struct {
	UserID   string `json:"userId"`
	UserUUID string `json:"userUUID"`
}

func (a *Account) ObjectID() primitive.ObjectID {
	p, err := primitive.ObjectIDFromHex(a.UserID)
	if err != nil {
		log.Fatalf("error convert the string to primitive.ObjectID: %+v", err)
	}
	return p
}
