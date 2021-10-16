package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/nECOnetic/storage/internal/service"
)

// StoreEcoPrediction ...
func (s *storage) StoreProfilerData(ctx context.Context, dataList []service.ProfilerData) error {
	// TODO: one transaction
	for _, data := range dataList {
		query, update := updateProfileData(data)

		opts := options.
			Update().
			SetUpsert(true)
		_, err := s.profilerDataCollection.UpdateOne(ctx, query, update, opts)
		if err != nil {
			return nil
		}
	}
	return nil
}
