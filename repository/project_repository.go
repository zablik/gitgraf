package repository

import (
	"context"
	"gitgraf/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProjectRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (*model.Project, error)
	Create(ctx context.Context, project *model.Project) (primitive.ObjectID, error)
	Update(ctx context.Context, project *model.Project) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}
