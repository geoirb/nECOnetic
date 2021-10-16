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

func ecoDataToMongo(srcList []service.EcoData) []interface{} {
	dst := make([]interface{}, 0, len(srcList))

	for _, src := range srcList {
		el := ecoData{
			Datatime:             src.Datatime,
			Measurement:          src.Measurement,
			PredictedMeasurement: src.PredictedMeasurement,
		}
		el.StationID, _ = primitive.ObjectIDFromHex(src.StationID)
		dst = append(dst, el)
	}
	return dst
}
