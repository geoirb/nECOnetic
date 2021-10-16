package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Storage mongo.
type storage struct {
	faceSearchCollection *mongo.Collection
}

// NewStorage ...
func NewStorage(
	ctx context.Context,
	connStr, databaseName, ecoCollectionName string,
) (*storage, error) {
	opts := options.Client().ApplyURI(connStr)
	connect, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}
	if err = connect.Ping(ctx, nil); err != nil {
		err = fmt.Errorf("error ping mongo storage %w", err)
	}

	collection := connect.Database(databaseName).Collection(ecoCollectionName)
	if _, err := collection.Indexes().CreateOne(ctx,
		mongo.IndexModel{
			Keys:    bson.M{"photo_hash": 1},
			Options: options.Index().SetUnique(true),
		}); err != nil {
		return nil, err
	}

	return &storage{
		faceSearchCollection: collection,
	}, err
}
