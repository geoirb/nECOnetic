package http

import (
	"github.com/gorilla/mux"

	"github.com/nECOnetic/data-service/internal/service"
)

var (
	prefix = "/api/v1/data-service"

	AddStationURI     = prefix + "/station"
	GetStationListURI = prefix + "/station"

	AddStationDataURI = prefix + "/station/data"
	GetStationDataURI = prefix + "/station/eco-data"
)

func Routing(r *mux.Router, svc service.Storage) {

}
