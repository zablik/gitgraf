package mongodb

import (
	"context"
	"gitgraf/model"
	"gitgraf/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type projectRepository struct {
	collection *mongo.Collection
}

func NewProjectRepository(collection *mongo.Collection) repository.ProjectRepository {
	return &projectRepository{collection: collection}
}

func (r *projectRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*model.Project, error) {
	filter := bson.M{"_id": id}
	project := &model.Project{}
	err := r.collection.FindOne(ctx, filter).Decode(project)
	return project, err
}

func (r *projectRepository) Create(ctx context.Context, project *model.Project) (primitive.ObjectID, error) {
	res, err := r.collection.InsertOne(ctx, project)
	if err != nil {
		return primitive.NilObjectID, err
	}
	project.ID = res.InsertedID.(primitive.ObjectID)
	return project.ID, nil
}

func (r *projectRepository) Update(ctx context.Context, project *model.Project) error {
	filter := bson.M{"_id": project.ID}
	update := bson.M{
		"$set": bson.M{
			"name":  project.Name,
			"title": project.Title,
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *projectRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}
