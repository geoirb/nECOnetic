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

func stationFilter(filter service.StationFilter) bson.M {
	return bson.M{
		"name": filter.Name,
	}
}
