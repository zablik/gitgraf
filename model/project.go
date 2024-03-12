package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name"`
	Title string             `bson:"title"`
}
