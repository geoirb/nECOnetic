package main

import (
	"context"
	"fmt"
	"log"
	"os"

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
		filePath:    "dataset/Академика Анохина 2020.xlsx",
	},
	{
		stationName: "Бутлерова",
		fileName:    "Бутлерова 2020.xlsx",
		filePath:    "dataset/Бутлерова 2020.xlsx",
	},
	{
		stationName: "Глебовская",
		fileName:    "Глебовская 2020 год.xlsx",
		filePath:    "dataset/Глебовская 2020 год.xlsx",
	},
	{
		stationName: "Коптевский",
		fileName:    "Коптевский бул. 2020 год.xlsx",
		filePath:    "dataset/Коптевский бул. 2020 год.xlsx",
	},
	{
		stationName: "Марьино",
		fileName:    "Марьино 2020.xlsx",
		filePath:    "dataset/Марьино 2020.xlsx",
	},
	{
		stationName: "Останкино",
		fileName:    "Останкино 0 2020 год.xlsx",
		filePath:    "dataset/Останкино 0 2020 год.xlsx",
	},
	{
		stationName: "Пролетарский",
		fileName:    "Пролетарский проспект 2020.xlsx",
		filePath:    "dataset/Пролетарский проспект 2020.xlsx",
	},
	{
		stationName: "Спиридоновка",
		fileName:    "Спиридоновка ул. 2020 год.xlsx",
		filePath:    "dataset/Спиридоновка ул. 2020 год.xlsx",
	},
	{
		stationName: "Туристская",
		fileName:    "Туристская 2020 год.xlsx",
		filePath:    "dataset/Туристская 2020 год.xlsx",
	},
	{
		stationName: "Шаболовка",
		fileName:    "Шаболовка 2020 год.xlsx",
		filePath:    "dataset/Шаболовка 2020.xlsx",
	},
}

func main() {
	f := mongo.Fabric{
		EcoDataCollectionName: "eco-data",
	}

	st, err := f.NewStorage(
		context.Background(),
		"mongodb://neconetic:neconetic@127.0.0.1:27017",
		"neconetic",
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
			Type:        "eco",
		}

		fmt.Println(svc.AddDataFromStation(context.Background(), data))
	}
}
