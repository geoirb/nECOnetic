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
	stationCollection      *mongo.Collection
	ecoDataCollection      *mongo.Collection
	profilerDataCollection *mongo.Collection
}

// NewStorage ...
func (f *Fabric) NewStorage(
	ctx context.Context,
	connStr, databaseName string,
) (*storage, error) {
	opts := options.Client().ApplyURI(connStr)
	connect, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}
	if err = connect.Ping(ctx, nil); err != nil {
		err = fmt.Errorf("error ping mongo storage %w", err)
		return nil, err
	}

	db := connect.Database(databaseName)

	s := &storage{}
	s.stationCollection = db.Collection(f.StationCollectionName)
	s.stationCollection.Indexes().CreateMany(
		ctx,
		[]mongo.IndexModel{
			{
				Keys: bson.M{
					"name": 1,
				},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys: bson.M{
					"lat": 1,
					"lon": 1,
				},
				Options: options.Index().SetUnique(true),
			},
		},
	)

	s.profilerDataCollection = db.Collection(f.ProfilerDataCollectionName)
	s.profilerDataCollection.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bson.M{
				"station_id": 1,
				"datatime":   1,
			},
			Options: options.Index().SetUnique(true),
		},
	)

	s.ecoDataCollection = db.Collection(f.EcoDataCollectionName)
	s.ecoDataCollection.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bson.M{
				"station_id": 1,
				"datatime":   1,
			},
			Options: options.Index().SetUnique(true),
		},
	)
	// TODO:
	// 1. station id must existing in station collection
	// 2. unique station id and datatime

	return s, nil
}
