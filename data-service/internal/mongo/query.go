package mongo

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/nECOnetic/data-service/internal/service"
)

func updateProfileData(src service.ProfilerData) (bson.M, bson.M) {
	filter := bson.M{
		"timestamp": src.Timestamp,
	}
	filter["station_id"], _ = primitive.ObjectIDFromHex(src.StationID)

	set := bson.M{
		"station_id":  filter["station_id"],
		"timestamp":   src.Timestamp,
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
		"timestamp": src.Timestamp,
	}
	filter["station_id"], _ = primitive.ObjectIDFromHex(src.StationID)

	set := bson.M{
		"station_id":            filter["station_id"],
		"timestamp":             src.Timestamp,
		"predicted_measurement": src.PredictedMeasurement,
	}

	if len(src.Measurement) != 0 {
		set["predicted_measurement"] = nil
		set["measurement"] = src.Measurement
	}

	return filter, bson.M{"$set": set}
}

func stationFilter(filter service.StationFilter) bson.M {
	f := make(bson.M)
	if filter.Name != nil {
		f["name"] = *filter.Name
	}
	return f
}

func ecoDataFilter(filter service.EcoDataFilter) bson.M {
	f := make(bson.M)

	if filter.StationID != nil {
		f["station_id"] = *filter.StationID
	}

	if filter.TimestampFrom != nil {
		f["timestamp"] = bson.M{
			"$gte": *filter.TimestampFrom,
		}
	}

	if filter.TimestampTo != nil {
		f["timestamp"] = bson.M{
			"$lte": *filter.TimestampTo,
		}
	}

	if len(filter.Measurements) != 0 {
		for _, field := range filter.Measurements {
			f["$or"] = []bson.M{
				{fmt.Sprintf("measurement.%s", field): bson.M{
					"$exists": true,
				},
				},
				{fmt.Sprintf("predicted_measurement.%s", field): bson.M{
					"$exists": true,
				},
				},
			}
		}
	}
	return f
}

func profilerDataFilter(filter service.ProfilerDataFilter) bson.M {
	f := bson.M{}

	if filter.StationID != nil {
		f["station_id"] = *filter.StationID
	}

	if filter.TimestampFrom != nil {
		f["timestamp"] = bson.M{
			"$gte": *filter.TimestampFrom,
		}
	}

	if filter.TimestampTo != nil {
		f["timestamp"] = bson.M{
			"$lte": *filter.TimestampTo,
		}
	}
	return f
}
