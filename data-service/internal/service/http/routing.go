package http

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/nECOnetic/data-service/internal/service"
)

var (
	prefix = "/api/v1/data-service"

	addStationURI     = prefix + "/station"
	getStationListURI = prefix + "/station"
	addStationDataURI = prefix + "/station/data"
	getEcoDataListURI = prefix + "/station/eco-data"
)

type buildResponseFunc func(payload interface{}, err error) ([]byte, error)

// Routing to svc.
func Routing(r *mux.Router, svc service.Storage, build buildResponseFunc) {
	r.Handle(addStationURI, addStationHandler(svc, build)).Methods(http.MethodPost)
	r.Handle(getStationListURI, getStationListHandler(svc, build)).Methods(http.MethodGet)

	r.Handle(addStationDataURI, addStationDataHandler(svc, build)).Methods(http.MethodPost)
	r.Handle(getEcoDataListURI, getEcoDataListHandler(svc, build)).Methods(http.MethodGet)
}
