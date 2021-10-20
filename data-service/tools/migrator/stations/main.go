package main

import (
	"context"
	"log"
	"time"

	"github.com/nECOnetic/data-service/internal/mongo"
	"github.com/nECOnetic/data-service/internal/service"
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
		Name: "Коптевский",
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
		Name: "Пролетарский",
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

func main() {
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

	for _, station := range stations {
		log.Println(st.StoreStation(context.Background(), station))
	}
}
