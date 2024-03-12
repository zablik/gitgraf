package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Commit struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty"`
	UserId    primitive.ObjectID   `bson:"user_id"`
	CreatedAt time.Time            `bson:"created_at"`
	Approvers []primitive.ObjectID `bson:"approvers"`
	Stats     Stats                `bson:"stats"`
}
