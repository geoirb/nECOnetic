package main

import (
	"context"
	"log"
	"os"
	"time"

	l "github.com/go-kit/log"
	"github.com/go-kit/log/level"

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
		filePath:    "../dataset/profiler/wind/03_метео_Останкино.xlsx",
	},
}

func main() {
	logger := l.NewJSONLogger(l.NewSyncWriter(os.Stdout))
	f := mongo.StorageFabric{
		StationCollectionName:      "station",
		EcoDataCollectionName:      "eco-data",
		ProfilerDataCollectionName: "profiler-data",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	st, err := f.NewStorage(
		ctx,
		"mongodb://localhost:27017/?readPreference=primary&ssl=false",
		"neconetic",
		7000,
	)
	if err != nil {
		log.Fatal(err)
	}

	svc := service.New(
		context.Background(),
		st,
		logger,
	)

	for _, src := range sources {
		file, err := os.Open(src.filePath)
		if err != nil {
			log.Fatal(err)
		}

		stationFilter := service.StationFilter{
			Name: &src.stationName,
		}

		stations, err := st.LoadStationList(ctx, stationFilter)
		if err != nil {
			log.Fatal(err)
		}

		dataList, err := svc.WindParser(ctx, stations[0].ID, src.fileName, file)
		if err != nil {
			log.Fatal(err)
		}

		start := time.Now()
		level.Debug(logger).Log("msg", "start store", "type", "eco")
		if err = st.StoreProfilerData(context.Background(), dataList); err != nil {
			log.Fatal(err)
		}
		level.Debug(logger).Log("msg", "finish store", "type", "eco", time.Since(start).Seconds())
	}
}
