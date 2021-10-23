package http

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nECOnetic/data-service/internal/service"
)

var (
	prefix = "/api/v1/data-service"

	addStationURI          = prefix + "/station"
	getStationListURI      = prefix + "/station"
	addStationDataURI      = prefix + "/station/data"
	getEcoDataListURI      = prefix + "/eco"
	getProfilerDataListURI = prefix + "/profiler"
	addPredictDataURI      = prefix + "/predict"
	// predictURI = prefix + "/station/profiler-data"
)

type bodyEncodeFunc func(payload interface{}, err error) ([]byte, error)

type svc interface {
	AddStation(ctx context.Context, in service.Station) (service.Station, error)
	AddDataFromStation(ctx context.Context, in service.StationData) error
	AddPredictedData(ctx context.Context, in []service.EcoData) error

	// Predict(ctx context.Context, in PredictFilter) error

	GetEcoDataList(ctx context.Context, in service.GetEcoData) ([]service.EcoData, error)
	GetProfilerDataList(ctx context.Context, in service.GetProfilerData) ([]service.ProfilerData, error)
	GetStationList(ctx context.Context) ([]service.Station, error)
}

// Routing to svc.
func Routing(r *mux.Router, svc svc, e bodyEncodeFunc) {
	r.Handle(addStationURI, addStationHandler(svc, e)).Methods(http.MethodPost)
	r.Handle(getStationListURI, getStationListHandler(svc, e)).Methods(http.MethodGet)

	r.Handle(addStationDataURI, addStationDataHandler(svc, e)).Methods(http.MethodPost)
	r.Handle(getEcoDataListURI, getEcoDataListHandler(svc, e)).Methods(http.MethodGet)
	r.Handle(getProfilerDataListURI, getProfilerDataListHandler(svc, e)).Methods(http.MethodGet)
	r.Handle(addPredictDataURI, addPredictDataHandler(svc, e)).Methods(http.MethodPost)

	// r.Handle(predictURI, predictHandler(svc, e)).Methods(http.MethodGet)
}
