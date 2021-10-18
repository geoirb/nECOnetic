package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/nECOnetic/data-service/internal/service"
)

// StoreProfilerData ...
func (s *storage) StoreProfilerData(ctx context.Context, dataList []service.ProfilerData) error {
	for _, data := range dataList {
		query, update := updateProfileData(data)

		opts := options.
			Update().
			SetUpsert(true)
		if _, err := s.profilerDataCollection.UpdateOne(ctx, query, update, opts); err != nil {
			return err
		}
	}
	return nil
}

// LoadProfilerData from storage.
func (s *storage) LoadProfilerData(ctx context.Context, filter service.ProfilerDataFilter) ([]service.ProfilerData, error) {
	f := profilerDataFilter(filter)

	cursor, err := s.profilerDataCollection.Find(ctx, f)
	// Check not found
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	data := make([]service.ProfilerData, 0, cursor.RemainingBatchLength())
	for cursor.Next(ctx) {
		var el profilerData
		if err = cursor.Decode(&el); err != nil {
			return nil, err
		}
		data = append(
			data,
			service.ProfilerData{
				StationID:          el.StationID.Hex(),
				Timestamp:          el.Timestamp,
				Temperature:        el.Temperature,
				OutsideTemperature: el.OutsideTemperature,
				WindDirection:      el.WindDirection,
				WindSpeed:          el.WindSpeed,
			},
		)
	}
	return data, err
}

// // StoreProfilerData ...
// func (s *storage) StoreProfilerData(ctx context.Context, dataList []service.ProfilerData) error {
// 	session, err := s.profilerDataCollection.Database().Client().StartSession()
// 	if err != nil {
// 		return err
// 	}
// 	defer session.EndSession(ctx)
// 	if err = session.StartTransaction(); err != nil {
// 		return err
// 	}

// 	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
// 		for _, data := range dataList {
// 			query, update := updateProfileData(data)

// 			opts := options.
// 				Update().
// 				SetUpsert(true)
// 			if _, err := s.profilerDataCollection.UpdateOne(sc, query, update, opts); err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	})
// 	return err
// }
