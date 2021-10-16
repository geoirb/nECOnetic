package mongo

import (
	"github.com/nECOnetic/storage/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
			Datatime:    src.Datatime,
			Measurement: src.Measurement,
		}
		el.StationID, _ = primitive.ObjectIDFromHex(src.StationID)
		dst = append(dst, el)
	}
	return dst
}

func ecoPredictionToMongo(srcList []service.EcoPrediction) []interface{} {
	dst := make([]interface{}, 0, len(srcList))
	for _, src := range srcList {
		el := ecoPrediction{
			ecoData{
				Datatime:    src.Datatime,
				Measurement: src.Measurement,
			},
		}
		el.StationID, _ = primitive.ObjectIDFromHex(src.StationID)
		dst = append(dst, el)
	}
	return dst
}
