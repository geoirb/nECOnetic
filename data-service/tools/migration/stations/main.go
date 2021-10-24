package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/kelseyhightower/envconfig"

	"github.com/nECOnetic/data-service/internal/mongo"
	"github.com/nECOnetic/data-service/internal/service"
)

var stations []service.Station = []service.Station{
	{
		Name: "Академика Анохина",
		Lat:  55.658163,
		Lon:  37.471434,
	},
	{
		Name: "Бутлерова",
		Lat:  55.649412,
		Lon:  37.535874,
	},
	{
		Name: "Глебовская",
		Lat:  55.811801,
		Lon:  37.71249,
	},
	{
		Name: "Коптевский",
		Lat:  55.833222,
		Lon:  37.525158,
	},
	{
		Name: "Марьино",
		Lat:  55.652695,
		Lon:  37.751502,
	},
	{
		Name: "Останкино",
		Lat:  55.821154,
		Lon:  37.612592,
	},
	{
		Name: "Пролетарский",
		Lat:  55.635129,
		Lon:  37.658684,
	},
	{
		Name: "Туристская",
		Lat:  55.856324,
		Lon:  37.426628,
	},
	{
		Name: "Спиридоновка",
		Lat:  55.759354,
		Lon:  37.595584,
	},
	{
		Name: "Шаболовка",
		Lat:  55.715698,
		Lon:  37.6052377,
	},
}

type configuration struct {
	StorageURI      string `envconfig:"STORAGE_URI" default:"mongodb://localhost:27017/?readPreference=primary&ssl=false"`
	StorageDatabase string `envconfig:"STORAGE_DATABASE" default:"neconetic"`

	StationCollectionName      string `envconfig:"STATION_COLLECTION_NAME" default:"station"`
	EcoDataCollectionName      string `envconfig:"ECO_DATA_COLLECTION_NAME" default:"eco-data"`
	ProfilerDataCollectionName string `envconfig:"PROFILER_DATA_COLLECTION_NAME" default:"profiler-data"`
}

func main() {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))

	var cfg configuration
	if err := envconfig.Process("", &cfg); err != nil {
		level.Error(logger).Log("msg", "configuration", "err", err)
		os.Exit(1)
	}

	level.Error(logger).Log("msg", "initialization", "cfg", cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	f := mongo.StorageFabric{
		StationCollectionName:      cfg.StationCollectionName,
		EcoDataCollectionName:      cfg.EcoDataCollectionName,
		ProfilerDataCollectionName: cfg.ProfilerDataCollectionName,
	}

	storage, err := f.NewStorage(
		ctx,
		cfg.StorageURI,
		cfg.StorageDatabase,
	)
	if err != nil {
		level.Error(logger).Log("msg", "init mongo", "err", err)
		os.Exit(1)
	}

	for _, station := range stations {
		fmt.Println(storage.StoreStation(context.Background(), station))
	}
}
