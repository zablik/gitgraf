package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Total struct {
	UserId  primitive.ObjectID `bson:"user_id"`
	RepoId  primitive.ObjectID `bson:"repo_id"`
	Commits int                `bson:"commits"`
	Reviews int                `bson:"reviews"`
	Stats   Stats              `bson:"stats"`
}
