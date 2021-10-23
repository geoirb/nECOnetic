package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	l "github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/nECOnetic/data-service/internal/mongo"
	"github.com/nECOnetic/data-service/internal/service"
)

var (
	stationName = "Останкино"
	srcDir      = "/home/geoirb/project/nECOnetic/dataset/profiler/temperature"
)

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
	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".txt") {
			file, err := os.Open(srcDir + "/" + f.Name())
			if err != nil {
				log.Fatal(err)
			}

			stationFilter := service.StationFilter{
				Name: &stationName,
			}

			stations, err := st.LoadStationList(ctx, stationFilter)
			if err != nil {
				log.Fatal(err)
			}

			dataList, err := svc.TemperatureParser(ctx, stations[0].ID, f.Name(), file)
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
}
