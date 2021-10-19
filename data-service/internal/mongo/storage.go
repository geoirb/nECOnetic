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
	hosts []string,
	username string,
	password string,
	databaseName string,
	transactionNumb int,
) (*storage, error) {

	opts := options.Client().
		SetAuth(
			options.Credential{
				Username: username,
				Password: password,
			},
		).
		SetHosts(hosts)

	connect, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}
	if err = connect.Ping(ctx, nil); err != nil {
		err = fmt.Errorf("error ping mongo storage %w", err)
		return nil, err
	}

	db := connect.Database(databaseName)

	s := &storage{
		transactionNumb: transactionNumb,
	}
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
				"timestamp":  1,
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
				"timestamp":  1,
			},
			Options: options.Index().SetUnique(true),
		},
	)
	// TODO:
	// 1. station id must existing in station collection
	// 2. unique station id and timestamp

	return s, nil
}
