package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/nECOnetic/data-service/internal/service"
)

func updateProfileData(src service.ProfilerData) (bson.M, bson.M) {
	filter := bson.M{
		"datatime": src.Datatime,
	}
	filter["_id"], _ = primitive.ObjectIDFromHex(src.StationID)

	set := bson.M{
		"_id":         filter["_id"],
		"datatime":    src.Datatime,
		"temperature": src.Temperature,
	}

	if src.OutsideTemperature != nil {
		set["outside_temperature"] = *src.OutsideTemperature
	}

	if src.WindDirection != nil {
		set["wind_direction"] = *src.WindDirection
	}

	if src.WindSpeed != nil {
		set["wind_speed"] = *src.WindSpeed
	}

	return filter, bson.M{"$set": set}
}

func updateEcoData(src service.EcoData) (bson.M, bson.M) {
	filter := bson.M{
		"datatime": src.Datatime,
	}
	filter["_id"], _ = primitive.ObjectIDFromHex(src.StationID)

	set := bson.M{
		"datatime":              src.Datatime,
		"predicted_measurement": src.PredictedMeasurement,
	}

	if len(src.Measurement) == 0 {
		set["predicted_measurement"] = nil
		set["measurement"] = src.Measurement
	}

	return filter, bson.M{"$set": set}
}

func stationFilter(filter service.StationFilter) bson.M {
	return bson.M{
		"name": filter.Name,
	}
}

func ecoDataFilter(filter service.EcoDataFilter) bson.M {
	f := bson.M{}

	if filter.StationID != nil {
		f["station_id"] = *filter.StationID
	}

	if filter.DatatimeFrom != nil {
		f["timestamp"] = bson.M{
			"$gte": *filter.DatatimeFrom,
		}
	}

	if filter.DatatimeTo != nil {
		f["timestamp"] = bson.M{
			"$lte": *filter.DatatimeTo,
		}
	}
	return f
}

func profilerDataFilter(filter service.ProfilerDataFilter) bson.M {
	f := bson.M{}

	if filter.StationID != nil {
		f["station_id"] = *filter.StationID
	}

	if filter.DatatimeFrom != nil {
		f["timestamp"] = bson.M{
			"$gte": *filter.DatatimeFrom,
		}
	}

	if filter.DatatimeTo != nil {
		f["timestamp"] = bson.M{
			"$lte": *filter.DatatimeTo,
		}
	}
	return f
}
