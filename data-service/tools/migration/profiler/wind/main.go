package main

import (
	"context"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/kelseyhightower/envconfig"

	"github.com/nECOnetic/data-service/internal/mongo"
	"github.com/nECOnetic/data-service/internal/service"
)

var sources []struct {
	stationName string
	fileName    string
	filePath    string
} = []struct {
	stationName string
	fileName    string
	filePath    string
}{
	{
		stationName: "Останкино",
		fileName:    "03_метео_Останкино.xlsx",
		filePath:    "/dataset/profiler/wind/03_метео_Останкино.xlsx",
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

	svc := service.New(
		context.Background(),
		storage,
		logger,
	)

	for _, src := range sources {
		file, err := os.Open(src.filePath)
		if err != nil {
			level.Error(logger).Log("open file", "file", "err", err)
			os.Exit(1)
		}

		stationFilter := service.StationFilter{
			Name: &src.stationName,
		}

		stations, err := storage.LoadStationList(context.Background(), stationFilter)
		if err != nil {
			level.Error(logger).Log("open file", "file", "err", err)
			os.Exit(1)
		}

		dataList, err := svc.WindParser(context.Background(), stations[0].ID, src.fileName, file)
		if err != nil {
			level.Error(logger).Log("open file", "file", "err", err)
			os.Exit(1)
		}

		start := time.Now()
		level.Debug(logger).Log("msg", "start store", "type", "wind")
		if err = storage.StoreProfilerData(context.Background(), dataList); err != nil {
			level.Error(logger).Log("msg", "store", "type", "wind", "err", err)
			os.Exit(1)
		}
		level.Debug(logger).Log("msg", "finish store", "type", "wind", "time", time.Since(start).Seconds())
	}
}
