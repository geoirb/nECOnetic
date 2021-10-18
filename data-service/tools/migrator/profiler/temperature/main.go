package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nECOnetic/data-service/internal/mongo"
	"github.com/nECOnetic/data-service/internal/service"
)

var (
	stationName = "Останкино"
	srcDir      = "../dataset/profiler"
)

func main() {
	f := mongo.Fabric{
		StationCollectionName:      "station",
		ProfilerDataCollectionName: "profiler-data",
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	st, err := f.NewStorage(
		ctx,
		[]string{
			"127.0.0.1:27017",
		},
		"neconetic",
		"neconetic",
		"neconetic",
	)
	if err != nil {
		log.Fatal(err)
	}

	svc := service.New(
		st,
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

			data := service.StationData{
				StationName: stationName,
				FileName:    f.Name(),
				File:        file,
				Type:        "temperature",
			}

			fmt.Println(svc.AddDataFromStation(context.Background(), data))
		}

	}
}
