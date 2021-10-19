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
	TimestampFrom *time.Time
	TimestampTo   *time.Time
	Measurements  []string
}

// ProfileData received from station.
type ProfilerData struct {
	StationID          string
	Timestamp          int64
	Temperature        map[string]float64
	OutsideTemperature *float64
	WindDirection      *int
	WindSpeed          *int
}

// ProfilerDataFilter ...
type ProfilerDataFilter struct {
	StationID     *string
	TimestampFrom *int64
	TimestampTo   *int64
}

type StationData struct {
	StationName string
	FileName    string
	File        io.Reader
	Type        string
}

type GetEcoData struct {
	StationID     *string
	TimestampFrom *time.Time
	TimestampTo   *time.Time
	Measurements  []string
}

type GetProfilerData struct {
	StationID     *string
	TimestampFrom *int64
	TimestampTo   *int64
}
