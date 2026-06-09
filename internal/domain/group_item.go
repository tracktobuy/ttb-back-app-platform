package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupItem struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	GroupId  primitive.ObjectID `bson:"groupId" json:"groupId"`
	ItemId   primitive.ObjectID `bson:"itemId" json:"itemId"`
	JoinedAt time.Time          `bson:"joinedAt" json:"joinedAt"`
}
