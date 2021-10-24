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

	transactionNumb int
}

// NewStorage ...
func (f *StorageFabric) NewStorage(
	ctx context.Context,
	uri string,
	databaseName string,
) (*storage, error) {

	opts := options.Client().ApplyURI(uri)

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
	_, err = s.stationCollection.Indexes().CreateMany(
		ctx,
		[]mongo.IndexModel{
			{
				Keys: bson.M{
					"name": 1,
				},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys: bson.D{
					{
						Key:   "lat",
						Value: 1,
					},
					{
						Key:   "lon",
						Value: 1,
					},
				},
				Options: options.Index().SetUnique(true),
			},
		},
	)
	if err != nil {
		return nil, err
	}

	s.profilerDataCollection = db.Collection(f.ProfilerDataCollectionName)
	_, err = s.profilerDataCollection.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bson.D{
				{
					Key:   "station_id",
					Value: 1,
				},
				{
					Key:   "timestamp",
					Value: 1,
				},
			},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		return nil, err
	}

	s.ecoDataCollection = db.Collection(f.EcoDataCollectionName)
	_, err = s.ecoDataCollection.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bson.D{
				{
					Key:   "station_id",
					Value: 1,
				},
				{
					Key:   "timestamp",
					Value: 1,
				},
			},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		return nil, err
	}

	// TODO:
	// 1. station id must existing in station collection
	// 2. unique station id and timestamp

	return s, err
}
