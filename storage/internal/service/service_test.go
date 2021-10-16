package service_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/nECOnetic/storage/internal/mongo"
	"github.com/nECOnetic/storage/internal/service"
)

var stations []service.Station = []service.Station{
	{
		Name: "Академика Анохина",
		Lat:  55.658163,
		Lon:  37.471434,
	},
	{
		Name: "Бутлерова",
		Lat:  55.658163,
		Lon:  37.471434,
	},
	{
		Name: "Глебовская",
		Lat:  55.811801,
		Lon:  37.71249,
	},
	{
		Name: "Коптевский бул",
		Lat:  55.833222,
		Lon:  37.525158,
	},
	{
		Name: "Марьино",
		Lat:  55.652695,
		Lon:  37.751502,
	},
	{
		Name: "Останкино",
		Lat:  55.821154,
		Lon:  37.612592,
	},
	{
		Name: "Пролетарский проспект",
		Lat:  55.635129,
		Lon:  37.658684,
	},
	{
		Name: "Туристская",
		Lat:  55.856324,
		Lon:  37.426628,
	},
	{
		Name: "Спиридоновка",
		Lat:  55.759354,
		Lon:  37.595584,
	},
	{
		Name: "Шаболовка",
		Lat:  55.715698,
		Lon:  37.6052377,
	},
}

func TestStoreStation(t *testing.T) {
	f := mongo.Fabric{
		StationCollectionName:       "station",
		EcoDataCollectionName:       "eco-data",
		ProfilerDataCollectionName:  "profiler",
		EcoPredictionCollectionName: "eco-prediction",
	}

	st, err := f.NewStorage(
		context.Background(),
		"mongodb://neconetic:neconetic@127.0.0.1:27017",
		"neconetic",
	)

	if err != nil {
		log.Fatal(err)
	}

	for _, station := range stations {
		log.Println(st.StoreStation(context.Background(), station))
	}
}

func TestStoreEcoData(t *testing.T) {
	f := mongo.Fabric{
		StationCollectionName:       "station",
		EcoDataCollectionName:       "eco-data",
		ProfilerDataCollectionName:  "profiler",
		EcoPredictionCollectionName: "eco-prediction",
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

	file, err := os.Open("/home/geoirb/project/hackaton/nECOnetic/dataset/Академика Анохина 2020.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	data := service.StationData{
		StationName: "Академика Анохина",
		FileName:    "Академика Анохина 2020.xlsx",
		File:        file,
		Type:        "eco",
	}

	fmt.Println(svc.AddDataFromStation(context.Background(), data))
}
