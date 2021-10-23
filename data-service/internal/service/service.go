package service

import (
	"context"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

var (
	ecoType         = "eco"
	temperatureType = "temperature"
	windType        = "wind"
)

type storage interface {
	StoreStation(ctx context.Context, st Station) (Station, error)
	StoreProfilerData(ctx context.Context, dataList []ProfilerData) error
	StoreEcoData(ctx context.Context, dataList []EcoData) error

	LoadStationList(ctx context.Context, f StationFilter) ([]Station, error)
	LoadEcoDataList(ctx context.Context, filter EcoDataFilter) ([]EcoData, error)
	LoadProfilerDataList(ctx context.Context, filter ProfilerDataFilter) ([]ProfilerData, error)
}

type predictClient interface {
	Predict(ctx context.Context, ecoData []EcoData, profilerData []ProfilerData) error
}

type service struct {
	ctx context.Context

	storage       storage
	predictClient predictClient

	logger log.Logger
}

// New returns Storage interface for work with storage.
func New(
	ctx context.Context,
	storage storage,
	logger log.Logger,
) *service {
	s := &service{
		ctx:     ctx,
		storage: storage,
		logger:  logger,
	}
	return s
}

// AddStation ...
func (s *service) AddStation(ctx context.Context, in Station) (Station, error) {
	return s.storage.StoreStation(ctx, in)
}

// AddDataFromStation parse data from station and put to storage.
func (s *service) AddDataFromStation(ctx context.Context, in StationData) error {
	defer in.File.Close()

	logger := log.WithPrefix(s.logger, "method", "AddDataFromStation")

	stationFilter := StationFilter{
		Name: &in.StationName,
	}

	stations, err := s.storage.LoadStationList(ctx, stationFilter)
	if err != nil {
		level.Error(logger).Log("msg", "load station from storage", "err", err)
		return err
	}

	var store func() error

	switch in.Type {
	case ecoType:
		dataList, err := s.EcoDataParser(ctx, stations[0].ID, in.FileName, in.File)
		if err != nil {
			return err
		}
		store = func() error {
			return s.storage.StoreEcoData(s.ctx, dataList)
		}
	case windType:
		dataList, err := s.WindParser(ctx, stations[0].ID, in.FileName, in.File)
		if err != nil {
			return err
		}
		store = func() error {
			return s.storage.StoreProfilerData(s.ctx, dataList)
		}
	case temperatureType:
		dataList, err := s.TemperatureParser(ctx, stations[0].ID, in.FileName, in.File)
		if err != nil {
			return err
		}
		store = func() error {
			return s.storage.StoreProfilerData(s.ctx, dataList)
		}
	default:
		return errUnknownType
	}

	go func() {
		start := time.Now()
		level.Debug(logger).Log("msg", "start store", "type", in.Type)
		if err = store(); err != nil {
			level.Error(logger).Log("msg", "store", "type", in.Type, "err", err)
			return
		}
		level.Debug(logger).Log("msg", "finish store", "type", in.Type, time.Since(start).Seconds())
	}()

	return nil
}

// AddPredictedData ...
func (s *service) AddPredictedData(ctx context.Context, in []EcoData) error {
	return s.storage.StoreEcoData(ctx, in)
}

// Predict measurements.
func (s *service) Predict(ctx context.Context, in PredictFilter) error {
	logger := log.WithPrefix(s.logger, "method", "Predict")

	stations, err := s.storage.LoadStationList(ctx, StationFilter{
		Name: in.StationName,
	})
	if err != nil {
		level.Error(logger).Log("msg", "load station fom storage", "err", err)
		return err
	}

	f := EcoDataFilter{
		TimestampFrom: &in.TimestampFrom,
		TimestampTill: &in.TimestampTill,
	}

	profilerData, err := s.storage.LoadProfilerDataList(ctx, ProfilerDataFilter{
		TimestampFrom: &in.TimestampFrom,
		TimestampTill: &in.TimestampTill,
	})
	if err != nil {
		level.Error(logger).Log("msg", "load profiler data fom storage", "err", err)
		return err
	}

	for _, station := range stations {
		f.StationID = &station.ID

		ecoData, err := s.storage.LoadEcoDataList(ctx, f)
		if err != nil {
			level.Error(logger).Log("msg", "load eco data fom storage", "err", err)
			return err
		}

		if err = s.predictClient.Predict(ctx, ecoData, profilerData); err != nil {
			level.Error(logger).Log("msg", "predict", "err", err)
			return err
		}
	}
	return nil
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
