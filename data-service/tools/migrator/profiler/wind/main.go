package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	l "github.com/go-kit/log"

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
	logger := l.NewJSONLogger(l.NewSyncWriter(os.Stdout))
	f := mongo.StorageFabric{
		StationCollectionName: "station",
		EcoDataCollectionName: "profiler-data",
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
		start := time.Now()
		fmt.Println(src, start)
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
		fmt.Println(time.Since(start).Minutes())
	}
	var a int
	fmt.Scan(&a)
}
