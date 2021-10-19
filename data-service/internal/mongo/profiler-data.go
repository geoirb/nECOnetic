package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/nECOnetic/data-service/internal/service"
)

// StoreProfilerDataTrx ...
func (s *storage) StoreProfilerDataTrx(ctx context.Context, dataList []service.ProfilerData) error {
	session, err := s.profilerDataCollection.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)
	if err = session.StartTransaction(); err != nil {
		return err
	}

	start := 0
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		for i, data := range dataList {
			if i%s.transactionNumb == 0 {
				start = i
				if err = session.StartTransaction(); err != nil {
					return err
				}
			}
			query, update := updateProfileData(data)

			opts := options.
				Update().
				SetUpsert(true)
			if _, err := s.profilerDataCollection.UpdateOne(sc, query, update, opts); err != nil {
				return err
			}

			if i == start+s.transactionNumb-1 || i == len(dataList)-1 {
				if err = session.CommitTransaction(sc); err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}

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

// LoadProfilerDataList ...
func (s *storage) LoadProfilerDataList(ctx context.Context, filter service.ProfilerDataFilter) ([]service.ProfilerData, error) {
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
