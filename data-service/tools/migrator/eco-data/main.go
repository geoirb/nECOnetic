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
		stationName: "Академика Анохина",
		fileName:    "Академика Анохина 2020.xlsx",
		filePath:    "../dataset/Академика Анохина 2020.xlsx",
	},
	// {
	// 	stationName: "Бутлерова",
	// 	fileName:    "Бутлерова 2020.xlsx",
	// 	filePath:    "../dataset/Бутлерова 2020.xlsx",
	// },
	// {
	// 	stationName: "Глебовская",
	// 	fileName:    "Глебовская 2020 год.xlsx",
	// 	filePath:    "../dataset/Глебовская 2020 год.xlsx",
	// },
	// {
	// 	stationName: "Коптевский",
	// 	fileName:    "Коптевский бул. 2020 год.xlsx",
	// 	filePath:    "../dataset/Коптевский бул. 2020 год.xlsx",
	// },
	// {
	// 	stationName: "Марьино",
	// 	fileName:    "Марьино 2020.xlsx",
	// 	filePath:    "../dataset/Марьино 2020.xlsx",
	// },
	// {
	// 	stationName: "Останкино",
	// 	fileName:    "Останкино 0 2020 год.xlsx",
	// 	filePath:    "../dataset/Останкино 0 2020 год.xlsx",
	// },
	// {
	// 	stationName: "Пролетарский",
	// 	fileName:    "Пролетарский проспект 2020.xlsx",
	// 	filePath:    "../dataset/Пролетарский проспект 2020.xlsx",
	// },
	// {
	// 	stationName: "Спиридоновка",
	// 	fileName:    "Спиридоновка ул. 2020 год.xlsx",
	// 	filePath:    "../dataset/Спиридоновка ул. 2020 год.xlsx",
	// },
	// {
	// 	stationName: "Туристская",
	// 	fileName:    "Туристская 2020 год.xlsx",
	// 	filePath:    "../dataset/Туристская 2020 год.xlsx",
	// },
	// {
	// 	stationName: "Шаболовка",
	// 	fileName:    "Шаболовка 2020 год.xlsx",
	// 	filePath:    "../dataset/Шаболовка 2020.xlsx",
	// },
}

func main() {
	f := mongo.Fabric{
		StationCollectionName: "station",
		EcoDataCollectionName: "eco-data",
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
		4000,
	)
	if err != nil {
		log.Fatal(err)
	}

	svc := service.New(
		st,
	)

	start := time.Now()
	for _, src := range sources {
		fmt.Println(src)
		file, err := os.Open(src.filePath)
		if err != nil {
			log.Fatal(err)
		}

		data := service.StationData{
			StationName: src.stationName,
			FileName:    src.fileName,
			File:        file,
			Type:        "eco",
		}

		if err = svc.AddDataFromStation(context.Background(), data); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println(time.Since(start).Minutes())
}
