package service

import (
	"io"
	"time"
)

// Station data.
type Station struct {
	ID   string
	Name string
	Lon  float32
	Lat  float32
}

type StationFilter struct {
	Name string
}

// EcoData received from station.
type EcoData struct {
	StationID   string
	Datatime    time.Time
	Measurement map[string]float64
}

// ProfileData received from station.
type ProfilerData struct {
	StationID          string
	Datatime           time.Time
	Temperature        map[int]float64
	OutsideTemperature *float64
	WindDirection      *int
	WindSpeed          *int
}

// EcoPrediction received from prediction module.
type EcoPrediction struct {
	EcoData
}

type StationData struct {
	StationName string
	FileName    string
	File        io.Reader
	Type        string
}
