package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserGroup struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	UserId   primitive.ObjectID `bson:"userId" json:"userId"`
	GroupId  primitive.ObjectID `bson:"groupId" json:"groupId"`
	JoinedAt time.Time          `bson:"joinedAt" json:"joinedAt"`
}
