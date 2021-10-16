package mongo

import (
	"context"

	"github.com/nECOnetic/storage/internal/service"
)

// StoreEcoPrediction ...
func (s *storage) StoreEcoPrediction(ctx context.Context, dataList []service.EcoPrediction) error {
	dataListDB := ecoPredictionToMongo(dataList)
	_, err := s.ecoPredictionCollection.InsertMany(ctx, dataListDB)
	return err
}
