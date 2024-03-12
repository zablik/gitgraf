package repository

import (
	"context"
	"gitgraf/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommitRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (*model.Commit, error)
	Create(ctx context.Context, commit *model.Commit) (primitive.ObjectID, error)
	Update(ctx context.Context, commit *model.Commit) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}
