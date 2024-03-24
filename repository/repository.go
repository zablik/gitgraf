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
	CreateMany(ctx context.Context, commits []*model.Commit) ([]primitive.ObjectID, error)
}

type ProjectRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (*model.Project, error)
	Create(ctx context.Context, project *model.Project) (primitive.ObjectID, error)
	Update(ctx context.Context, project *model.Project) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type UserRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (*model.User, error)
	Create(ctx context.Context, user *model.User) (primitive.ObjectID, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	CreateMany(ctx context.Context, commits []*model.User) ([]primitive.ObjectID, error)
	GetAll(ctx context.Context) ([]*model.User, error)
}
