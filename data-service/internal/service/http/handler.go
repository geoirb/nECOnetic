package http

import (
	"net/http"

	"github.com/nECOnetic/data-service/internal/service"
)

type addStationServer struct {
	svc       service.Storage
	transport *addStationTransport
}

// ServeHTTP for adding station at system.
func (s *addStationServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	addingStation, err := s.transport.DecodeRequest(r)

	var addedStation service.Station
	if err == nil {
		addedStation, err = s.svc.AddStation(r.Context(), addingStation)
	}

	s.transport.EncodeResponse(w, addedStation, err)
}

func addStationHandler(svc service.Storage, build buildResponseFunc) http.Handler {
	return &addStationServer{
		svc:       svc,
		transport: newAddStationTransport(build),
	}
}

type getStationListServer struct {
	svc       service.Storage
	transport *getStationListTransport
}

// ServeHTTP for getting station list.
func (s *getStationListServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stationList, err := s.svc.GetStationList(r.Context())

	s.transport.EncodeResponse(w, stationList, err)
}

func getStationListHandler(svc service.Storage, build buildResponseFunc) http.Handler {
	return &getStationListServer{
		svc:       svc,
		transport: newGetStationListTransport(build),
	}
}

type addStationDataServer struct {
	svc       service.Storage
	transport *addStationDataTransport
}

func (s *addStationDataServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stationData, err := s.transport.DecodeRequest(r)

	if err == nil {
		err = s.svc.AddDataFromStation(r.Context(), stationData)
	}
	s.transport.EncodeResponse(w, err)
}

func addStationDataHandler(svc service.Storage, build buildResponseFunc) http.Handler {
	return &addStationDataServer{
		svc:       svc,
		transport: newAddStationDataTransport(build),
	}
}

type getEcoDataListServer struct {
	svc       service.Storage
	transport *getEcoDataListTransport
}

func (s *getEcoDataListServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filter, err := s.transport.DecodeRequest(r)

	var data []service.EcoData
	if err == nil {
		data, err = s.svc.GetEcoDataList(r.Context(), filter)
	}

	s.transport.EncodeResponse(w, data, err)
}

func getEcoDataListHandler(svc service.Storage, build buildResponseFunc) http.Handler {
	return &getEcoDataListServer{
		svc:       svc,
		transport: newGetEcoDataListTransport(build),
	}
}
