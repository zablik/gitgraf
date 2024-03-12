package repository

import (
	"context"
	"gitgraf/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (*model.User, error)
	Create(ctx context.Context, user *model.User) (primitive.ObjectID, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}
