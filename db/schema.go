package db

import (
	"context"
	"fmt"
	"gitgraf/config"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupMongo(ctx context.Context, cfg *config.Config) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DB.Host+":"+cfg.DB.Port))

	if err != nil {
		log.Fatal(err)
	}

	// users
	usersCollection := client.Database(cfg.DB.Name).Collection("users")
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	indexName, err := usersCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created index:", indexName)

	// commits
	commitsCollection := client.Database(cfg.DB.Name).Collection("commits")
	indexModel = mongo.IndexModel{
		Keys:    bson.D{{Key: "hash", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	indexName, err = commitsCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created index:", indexName)

	return client, nil
}
