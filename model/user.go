package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Email         string             `bson:"email"`
	Name          string             `bson:"name"`
	AltEmails     []string           `bson:"alt_emails"`
	AltNames      []string           `bson:"alt_names"`
	FirstActiveAt time.Time          `bson:"first_active_at"`
	LastActiveAt  time.Time          `bson:"last_active_at"`
}
