package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/nECOnetic/data-service/internal/service"
)

// StoreEcoData in storage.
func (s *storage) StoreEcoData(ctx context.Context, dataList []service.EcoData) error {
	session, err := s.ecoDataCollection.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	start := 0
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		for i, data := range dataList {
			if i%s.operationInTransaction == 0 {
				start = i
				if err = session.StartTransaction(); err != nil {
					return err
				}
			}

			query, update := updateEcoData(data)
			opts := options.
				Update().
				SetUpsert(true)
			if _, err := s.ecoDataCollection.UpdateOne(sc, query, update, opts); err != nil {
				return err
			}

			if i == start+s.operationInTransaction-1 || i == len(dataList)-1 {
				if err = session.CommitTransaction(sc); err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}

// LoadEcoData from storage.
func (s *storage) LoadEcoData(ctx context.Context, filter service.EcoDataFilter) ([]service.EcoData, error) {
	f := ecoDataFilter(filter)

	cursor, err := s.ecoDataCollection.Find(ctx, f)
	// Check not found
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	data := make([]service.EcoData, 0, cursor.RemainingBatchLength())
	for cursor.Next(ctx) {
		var el ecoData
		if err = cursor.Decode(&el); err != nil {
			return nil, err
		}
		data = append(
			data,
			service.EcoData{
				StationID:            el.StationID.Hex(),
				Timestamp:            el.Timestamp,
				Measurement:          el.Measurement,
				PredictedMeasurement: el.PredictedMeasurement,
			},
		)
	}
	return data, err
}
