package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	l "github.com/go-kit/log"

	"github.com/nECOnetic/data-service/internal/mongo"
	"github.com/nECOnetic/data-service/internal/service"
)

var (
	stationName = "Останкино"
	srcDir      = "../dataset/profiler"
)

func main() {
	logger := l.NewJSONLogger(l.NewSyncWriter(os.Stdout))
	f := mongo.StorageFabric{
		StationCollectionName:      "station",
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
		st,
		logger,
	)
	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".txt") {
			start := time.Now()
			fmt.Println(f, start)
			file, err := os.Open(srcDir + "/" + f.Name())
			if err != nil {
				log.Fatal(err)
			}

			data := service.StationData{
				StationName: stationName,
				FileName:    f.Name(),
				File:        file,
				Type:        "temperature",
			}

			fmt.Println(svc.AddDataFromStation(context.Background(), data))
			fmt.Println(time.Since(start).Minutes())
		}

	}
}
