package mongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type station struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
	Lat  float64            `bson:"lat"`
	Lon  float64            `bson:"lon"`
}

type ecoData struct {
	StationID            primitive.ObjectID `bson:"station_id"`
	Timestamp            int64              `bson:"timestamp"`
	Measurement          map[string]float64 `bson:"measurement"`
	PredictedMeasurement map[string]float64 `bson:"predicted_measurement"`
}

type profilerData struct {
	StationID          primitive.ObjectID `bson:"station_id"`
	Timestamp          int64              `bson:"timestamp"`
	Temperature        map[string]float64 `bson:"temperature"`
	OutsideTemperature *float64           `bson:"outside_temperature,omitempty"`
	WindDirection      *int               `bson:"wind_direction,omitempty"`
	WindSpeed          *int               `bson:"wind_speed"`
}
