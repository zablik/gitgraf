package mongodb

import (
	"context"
	"gitgraf/model"
	"gitgraf/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type commitRepository struct {
	collection *mongo.Collection
}

func NewCommitRepository(client *mongo.Client, dbName, collectionName string) repository.CommitRepository {
	return &commitRepository{
		collection: client.Database(dbName).Collection(collectionName),
	}
}

func (r *commitRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*model.Commit, error) {
	filter := bson.M{"_id": id}
	commit := &model.Commit{}
	err := r.collection.FindOne(ctx, filter).Decode(commit)
	return commit, err
}
func (r *commitRepository) Create(ctx context.Context, commit *model.Commit) (primitive.ObjectID, error) {
	res, err := r.collection.InsertOne(ctx, commit)
	if err != nil {
		return primitive.NilObjectID, err
	}
	commit.ID = res.InsertedID.(primitive.ObjectID)
	return commit.ID, nil
}

func (r *commitRepository) Update(ctx context.Context, commit *model.Commit) error {
	filter := bson.M{"_id": commit.ID}
	update := bson.M{
		"$set": bson.M{
			"user_id":    commit.UserId,
			"created_at": commit.CreatedAt,
			"approver_ids":  commit.ApproverIds,
			"stats": bson.M{
				"lines_added":    commit.Stats.LinesAdded,
				"lines_deleted":  commit.Stats.LinesDeleted,
				"files_added":    commit.Stats.FilesAdded,
				"files_deleted":  commit.Stats.FilesDeleted,
				"files_modified": commit.Stats.FilesModified,
			},
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *commitRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

func (r *commitRepository) CreateMany(ctx context.Context, commits []*model.Commit) ([]primitive.ObjectID, error) {
	documents := make([]interface{}, len(commits))
	for i, commit := range commits {
		documents[i] = commit
	}

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