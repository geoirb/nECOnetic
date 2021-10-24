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
		stationName: "Академика Анохина",
		fileName:    "Академика Анохина 2020.xlsx",
		filePath:    "/dataset/Академика Анохина 2020.xlsx",
	},
	{
		stationName: "Бутлерова",
		fileName:    "Бутлерова 2020.xlsx",
		filePath:    "/dataset/Бутлерова 2020.xlsx",
	},
	{
		stationName: "Глебовская",
		fileName:    "Глебовская 2020 год.xlsx",
		filePath:    "/dataset/Глебовская 2020 год.xlsx",
	},
	{
		stationName: "Коптевский",
		fileName:    "Коптевский бул. 2020 год.xlsx",
		filePath:    "/dataset/Коптевский бул. 2020 год.xlsx",
	},
	{
		stationName: "Марьино",
		fileName:    "Марьино 2020.xlsx",
		filePath:    "/dataset/Марьино 2020.xlsx",
	},
	{
		stationName: "Останкино",
		fileName:    "Останкино 0 2020 год.xlsx",
		filePath:    "/dataset/Останкино 0 2018 год.xlsx",
	},
	{
		stationName: "Останкино",
		fileName:    "Останкино 0 2020 год.xlsx",
		filePath:    "/dataset/Останкино 0 2019 год.xlsx",
	},
	{
		stationName: "Останкино",
		fileName:    "Останкино 0 2020 год.xlsx",
		filePath:    "/dataset/Останкино 0 2020 год.xlsx",
	},
	{
		stationName: "Пролетарский",
		fileName:    "Пролетарский проспект 2020.xlsx",
		filePath:    "/dataset/Пролетарский проспект 2020.xlsx",
	},
	{
		stationName: "Спиридоновка",
		fileName:    "Спиридоновка ул. 2020 год.xlsx",
		filePath:    "/dataset/Спиридоновка ул. 2020 год.xlsx",
	},
	{
		stationName: "Туристская",
		fileName:    "Туристская 2020 год.xlsx",
		filePath:    "/dataset/Туристская 2020 год.xlsx",
	},
	{
		stationName: "Шаболовка",
		fileName:    "Шаболовка 2020 год.xlsx",
		filePath:    "/dataset/Шаболовка 2020.xlsx",
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
			level.Error(logger).Log("msg", "open file", "err", err)
			os.Exit(1)
		}

		stationFilter := service.StationFilter{
			Name: &src.stationName,
		}

		stations, err := storage.LoadStationList(context.Background(), stationFilter)
		if err != nil {
			level.Error(logger).Log("msg", "load station", "err", err)
			os.Exit(1)
		}

		dataList, err := svc.EcoDataParser(context.Background(), stations[0].ID, src.fileName, file)
		if err != nil {
			level.Error(logger).Log("msg", "parse data", "err", err)
			os.Exit(1)
		}

		start := time.Now()
		level.Debug(logger).Log("msg", "start store", "type", "eco", "n", len(dataList))
		if err = storage.StoreEcoData(context.Background(), dataList); err != nil {
			level.Error(logger).Log("msg", "store", "type", "eco", "err", err)
			os.Exit(1)
		}
		level.Debug(logger).Log("msg", "finish store", "type", "eco", "time", time.Since(start).Seconds())
	}
}
