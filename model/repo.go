package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repo struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	RepoID primitive.ObjectID `bson:"repo_id"`
	Name   string             `bson:"name"`
	Title  string             `bson:"title"`
}
