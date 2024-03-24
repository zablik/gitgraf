package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Commit struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty"`
	Hash        string               `bson:"hash"`
	UserId      primitive.ObjectID   `bson:"user_id"`
	CreatedAt   time.Time            `bson:"created_at"`
	ApproverIds []primitive.ObjectID `bson:"approver_ids"`
	Stats       Stats                `bson:"stats"`
}
