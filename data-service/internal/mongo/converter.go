package mongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/nECOnetic/data-service/internal/service"
)

func stationToMongo(src service.Station) station {
	return station{
		ID:   primitive.NewObjectID(),
		Name: src.Name,
		Lon:  src.Lon,
		Lat:  src.Lat,
	}
}


