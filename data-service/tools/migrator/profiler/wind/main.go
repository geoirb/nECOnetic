package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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
		fileName:    "03_метео_Останкино.xls",
		filePath:    "../dataset/profiler/03_метео_Останкино.xls",
	},
}

func main() {
	f := mongo.Fabric{
		StationCollectionName: "station",
		EcoDataCollectionName: "profiler-data",
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
		7000,
	)
	if err != nil {
		log.Fatal(err)
	}

	svc := service.New(
		st,
	)

	for _, src := range sources {
		file, err := os.Open(src.filePath)
		if err != nil {
			log.Fatal(err)
		}

		data := service.StationData{
			StationName: src.stationName,
			FileName:    src.fileName,
			File:        file,
			Type:        "wind",
		}

		fmt.Println(svc.AddDataFromStation(context.Background(), data))
	}
}
