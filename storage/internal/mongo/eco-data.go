package mongo

import (
	"context"

	"github.com/nECOnetic/storage/internal/service"
)

// StoreEcoData ...
func (s *storage) StoreEcoData(ctx context.Context, data []service.EcoData) error {
	dataDB := ecoDataToMongo(data)
	_, err := s.ecoDataCollection.InsertMany(ctx, dataDB)
	return err
}
