package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type station struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
	Lon  float32            `bson:"lon"`
	Lat  float32            `bson:"lat"`
}

type ecoData struct {
	StationID   primitive.ObjectID `bson:"station_id"`
	Datatime    time.Time          `bson:"datatime"`
	Measurement map[string]float64 `bson:"measurement"`
}

type profilerData struct {
	StationID          primitive.ObjectID `bson:"station_id"`
	Datatime           time.Time          `bson:"datatime"`
	Temperature        map[int]float64    `bson:"temperature"`
	OutsideTemperature *float64           `bson:"outside_temperature,omitempty"`
	WindDirection      *int               `bson:"wind_direction,omitempty"`
	WindSpeed          *int               `bson:"wind_speed"`
}

type ecoPrediction struct {
	ecoData
}
