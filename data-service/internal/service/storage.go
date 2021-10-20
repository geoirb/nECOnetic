package service

import (
	"context"
)

type Storage interface {
	AddStation(ctx context.Context, in Station) (Station, error)
	AddDataFromStation(ctx context.Context, in StationData) error
	AddPredictedData(ctx context.Context, in []EcoData) error

	GetEcoDataList(ctx context.Context, in GetEcoData) ([]EcoData, error)
	GetProfilerDataList(ctx context.Context, in GetProfilerData) ([]ProfilerData, error)
	GetStationList(ctx context.Context) ([]Station, error)
}
