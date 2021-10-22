package http

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/nECOnetic/data-service/internal/service"
)

var (
	prefix = "/api/v1/data-service"

	addStationURI          = prefix + "/station"
	getStationListURI      = prefix + "/station"
	addStationDataURI      = prefix + "/station/data"
	getEcoDataListURI      = prefix + "/station/eco-data"
	getProfilerDataListURI = prefix + "/station/profiler-data"

	predictURI = prefix + "/station/profiler-data"
)

type bodyEncodeFunc func(payload interface{}, err error) ([]byte, error)

// Routing to svc.
func Routing(r *mux.Router, svc service.Storage, e bodyEncodeFunc) {
	r.Handle(addStationURI, addStationHandler(svc, e)).Methods(http.MethodPost)
	r.Handle(getStationListURI, getStationListHandler(svc, e)).Methods(http.MethodGet)

	r.Handle(addStationDataURI, addStationDataHandler(svc, e)).Methods(http.MethodPost)
	r.Handle(getEcoDataListURI, getEcoDataListHandler(svc, e)).Methods(http.MethodGet)
	r.Handle(getProfilerDataListURI, getProfilerDataListHandler(svc, e)).Methods(http.MethodGet)

	r.Handle(predictURI, predictHandler(svc, e)).Methods(http.MethodGet)
}
