package db

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database
var userCollection string = "user"
var boardCollection string = "board"

func getDbContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func InitDb() error {
	var client *mongo.Client
	var err error

	var opts *options.ClientOptions = options.Client()
	opts.SetTimeout(5 * time.Second)
	opts.ApplyURI(os.Getenv(("MONGODB_URI")))

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	ctx, cancel := getDbContext()
	defer cancel()

	client, err = mongo.Connect(ctx, opts)
	if err != nil {
		return err
	}

	db = client.Database(os.Getenv("MONGODB_DBNAME"))
	return nil
}

func RunConfig() (string, error) {
	ctx, cancel := getDbContext()
	defer cancel()
	var unique bool = true

	errString, err := db.Collection(userCollection).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"email": 1,
		}, Options: &options.IndexOptions{
			Unique: &unique,
		},
	})

	if err != nil {
		return errString, err
	}

	return "", nil
}
