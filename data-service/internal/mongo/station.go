package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/nECOnetic/data-service/internal/service"
)

// TODO:
// delete station
// update station

// StoreStation ...
func (s *storage) StoreStation(ctx context.Context, st service.Station) (service.Station, error) {
	stDB := stationToMongo(st)
	res, err := s.stationCollection.InsertOne(ctx, stDB)
	if err != nil {
		return st, err
	}

	if res != nil {
		st.ID = res.InsertedID.(primitive.ObjectID).Hex()
	}
	return st, err
}

// LoadStationList ...
func (s *storage) LoadStationList(ctx context.Context, filter service.StationFilter) ([]service.Station, error) {
	f := stationFilter(filter)

	cursor, err := s.stationCollection.Find(ctx, f)
	// Check not found
	if err != nil {
		return nil, err
	}
	if cursor.RemainingBatchLength() == 0 {
		return nil, errStationNotFound
	}
	defer cursor.Close(ctx)

	stations := make([]service.Station, 0, cursor.RemainingBatchLength())
	for cursor.Next(ctx) {
		var el station
		if err = cursor.Decode(&el); err != nil {
			return nil, err
		}
		stations = append(
			stations,
			service.Station{
				ID:   el.ID.Hex(),
				Name: el.Name,
				Lon:  el.Lon,
				Lat:  el.Lat,
			},
		)
	}

	return stations, err
}
