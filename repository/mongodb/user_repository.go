package mongodb

import (
	"context"
	"gitgraf/model"
	"gitgraf/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) repository.UserRepository {
	return &userRepository{collection: collection}
}

func (r *userRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	filter := bson.M{"_id": id}
	user := &model.User{}
	err := r.collection.FindOne(ctx, filter).Decode(user)
	return user, err
}

func (r *userRepository) Create(ctx context.Context, user *model.User) (primitive.ObjectID, error) {
	res, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return primitive.NilObjectID, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user.ID, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	filter := bson.M{"_id": user.ID}
	update := bson.M{
		"$set": bson.M{
			"name": user.Name,
			
			// Add other user fields here...
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}
