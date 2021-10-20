package service

import (
	"context"
	"io"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type storage interface {
	StoreStation(ctx context.Context, st Station) (Station, error)
	StoreProfilerData(ctx context.Context, dataList []ProfilerData) error
	StoreEcoData(ctx context.Context, dataList []EcoData) error

	LoadStationList(ctx context.Context, f StationFilter) ([]Station, error)
	LoadEcoDataList(ctx context.Context, filter EcoDataFilter) ([]EcoData, error)
	LoadProfilerDataList(ctx context.Context, filter ProfilerDataFilter) ([]ProfilerData, error)
}


type service struct {
	ctx         context.Context
	dataHandler map[string]func(context.Context, string, string, io.Reader) (err error)

	storage storage

	logger log.Logger
}

// New returns Storage interface for work with storage.
func New(
	ctx context.Context,
	storage storage,
	logger log.Logger,
) Storage {
	s := &service{
		ctx:     ctx,
		storage: storage,
		logger:  logger,
	}

	s.dataHandler = map[string]func(context.Context, string, string, io.Reader) (err error){
		"eco":         s.ecoDataHandler,
		"wind":        s.windHandler,
		"temperature": s.temperatureHandler,
	}

	return s
}

// AddStation ...
func (s *service) AddStation(ctx context.Context, in Station) (Station, error) {
	return s.storage.StoreStation(ctx, in)
}

// AddDataFromStation parse data from station and put to storage.
func (s *service) AddDataFromStation(ctx context.Context, in StationData) error {
	logger := log.WithPrefix(s.logger, "method", "AddDataFromStation")

	stationFilter := StationFilter{
		Name: &in.StationName,
	}

	stations, err := s.storage.LoadStationList(ctx, stationFilter)
	if err != nil {
		level.Error(logger).Log("msg", "load station from storage", "err", err)
		return err
	}

	h, isExist := s.dataHandler[in.Type]
	if !isExist {
		return errUnknownType
	}
	defer in.File.Close()

	return h(ctx, stations[0].ID, in.FileName, in.File)
}

// AddPredictedData ...
func (s *service) AddPredictedData(ctx context.Context, in []EcoData) error {
	return s.storage.StoreEcoData(ctx, in)
}

// GetStationList ...
func (s *service) GetStationList(ctx context.Context) ([]Station, error) {
	return s.storage.LoadStationList(ctx, StationFilter{})
}

// GetEcoDataList ...
func (s *service) GetEcoDataList(ctx context.Context, in GetEcoData) ([]EcoData, error) {
	logger := log.WithPrefix(s.logger, "method", "GetEcoDataList")

	f := EcoDataFilter{

		TimestampFrom: in.TimestampFrom,
		TimestampTill: in.TimestampTill,
		Measurements:  in.Measurements,
	}

	if in.StationName != nil {
		stations, err := s.storage.LoadStationList(ctx, StationFilter{
			Name: in.StationName,
		})
		if err != nil {
			level.Error(logger).Log("msg", "load station fom storage", "err", err)
			return nil, err
		}

		f.StationID = &stations[0].ID
	}
	return s.storage.LoadEcoDataList(ctx, f)
}

// GetProfilerDataList ...
func (s *service) GetProfilerDataList(ctx context.Context, in GetProfilerData) ([]ProfilerData, error) {
	logger := log.WithPrefix(s.logger, "method", "GetEcoDataList")

	f := ProfilerDataFilter{
		TimestampFrom: in.TimestampFrom,
		TimestampTill: in.TimestampTill,
	}

	if in.StationName != nil {
		stations, err := s.storage.LoadStationList(ctx, StationFilter{
			Name: in.StationName,
		})
		if err != nil {
			level.Error(logger).Log("msg", "load station fom storage", "err", err)
			return nil, err
		}

		f.StationID = &stations[0].ID
	}
	return s.storage.LoadProfilerDataList(ctx, f)
}
