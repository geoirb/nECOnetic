package mongo

import (
	"context"

	"github.com/nECOnetic/data-service/internal/service"
)

// StoreEcoData ...
func (s *storage) StoreEcoData(ctx context.Context, data []service.EcoData) error {
	// TODO:
	// upsert add
	dataDB := ecoDataToMongo(data)
	_, err := s.ecoDataCollection.InsertMany(ctx, dataDB)
	return err
}
