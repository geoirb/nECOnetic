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

// StationFilter ...
type StationFilter struct {
	Name string
}

// EcoData received from station.
type EcoData struct {
	StationID            string
	Datatime             int64
	Measurement          map[string]float64
	PredictedMeasurement map[string]float64
}

// EcoDataFilter ...
type EcoDataFilter struct {
	StationID    *string
	DatatimeFrom *time.Time
	DatatimeTo   *time.Time
}

// ProfileData received from station.
type ProfilerData struct {
	StationID          string
	Datatime           int64
	Temperature        map[int]float64
	OutsideTemperature *float64
	WindDirection      *int
	WindSpeed          *int
}

// ProfilerDataFilter ...
type ProfilerDataFilter struct {
	StationID    *string
	DatatimeFrom *int64
	DatatimeTo   *int64
}

type StationData struct {
	StationName string
	FileName    string
	File        io.Reader
	Type        string
}
