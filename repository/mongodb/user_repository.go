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

func NewUserRepository(client *mongo.Client, dbName, collectionName string) repository.UserRepository {
	return &userRepository{
		collection: client.Database(dbName).Collection(collectionName),
	}
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

func (r *userRepository) CreateMany(ctx context.Context, users []*model.User) ([]primitive.ObjectID, error) {
	documents := convertToInterfaceSlice(users)

	res, err := r.collection.InsertMany(ctx, documents)
	if err != nil {
		return nil, err
	}

	ids := make([]primitive.ObjectID, len(res.InsertedIDs))
	for i, id := range res.InsertedIDs {
		ids[i] = id.(primitive.ObjectID)
	}

	return ids, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]*model.User, error) {
	filter := bson.M{} // empty filter to get all documents
	users := []*model.User{}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		user := &model.User{}
		err := cursor.Decode(user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func convertToInterfaceSlice(slice []*model.User) []interface{} {
	var s []interface{}
	for _, item := range slice {
		s = append(s, item)
	}
	return s
}
