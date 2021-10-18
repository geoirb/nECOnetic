package service

import (
	"context"
	"errors"
	"io"
)

type storage interface {
	StoreStation(ctx context.Context, st Station) (Station, error)
	StoreProfilerData(ctx context.Context, dataList []ProfilerData) error
	StoreEcoData(ctx context.Context, dataList []EcoData) error

	LoadStation(ctx context.Context, f StationFilter) ([]Station, error)
}

// TODO:
// Logging

type service struct {
	storage storage

	dataHandler map[string]func(context.Context, string, io.Reader) (err error)
}

// New returns Storage interface for work with storage.
func New(
	storage storage,
) Storage {
	s := &service{
		storage: storage,
	}

	s.dataHandler = map[string]func(context.Context, string, io.Reader) (err error){
		"eco":         s.ecoDataHandler,
		"wind":        s.windHandler,
		"temperature": s.temperatureHandler,
	}

	return s
}

// AddDataFromStation to storage.
func (s *service) AddDataFromStation(ctx context.Context, in StationData) error {
	stationFilter := StationFilter{
		Name: in.StationName,
	}

	stations, err := s.storage.LoadStation(ctx, stationFilter)
	if err != nil {
		return err
	}

	h, isExist := s.dataHandler[in.Type]
	if !isExist {
		return errors.New("unknown data type")
	}

	return h(ctx, stations[0].ID, in.File)
}

// AddPredictedData to storage.
func (s *service) AddPredictedData(ctx context.Context, in []EcoData) error {
	return s.storage.StoreEcoData(ctx, in)
}
