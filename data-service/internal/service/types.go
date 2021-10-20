package service

import (
	"io"
)

// Station data.
type Station struct {
	ID   string
	Name string
	Lat  float64
	Lon  float64
}

// StationFilter ...
type StationFilter struct {
	Name *string
}

// EcoData received from station.
type EcoData struct {
	StationID            string
	Timestamp            int64
	Measurement          map[string]float64
	PredictedMeasurement map[string]float64
}

// EcoDataFilter ...
type EcoDataFilter struct {
	StationID     *string
	TimestampFrom *int64
	TimestampTill *int64
	Measurements  []string
}

// ProfileData received from station.
type ProfilerData struct {
	StationID          string
	Timestamp          int64
	Temperature        map[string]float64
	OutsideTemperature *float64
	WindDirection      *int
	WindSpeed          *float64
}

// ProfilerDataFilter ...
type ProfilerDataFilter struct {
	StationID     *string
	TimestampFrom *int64
	TimestampTill *int64
}

type StationData struct {
	StationName string
	FileName    string
	File        io.ReadCloser
	Type        string
}

type GetEcoData struct {
	StationName   *string
	TimestampFrom *int64
	TimestampTill *int64
	Measurements  []string
}

type GetProfilerData struct {
	StationID     *string
	TimestampFrom *int64
	TimestampTill *int64
}
