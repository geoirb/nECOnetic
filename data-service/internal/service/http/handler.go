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

func addStationHandler(svc service.Storage, be bodyEncodeFunc) http.Handler {
	return &addStationServer{
		svc:       svc,
		transport: newAddStationTransport(be),
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

func getStationListHandler(svc service.Storage, be bodyEncodeFunc) http.Handler {
	return &getStationListServer{
		svc:       svc,
		transport: newGetStationListTransport(be),
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

func addStationDataHandler(svc service.Storage, be bodyEncodeFunc) http.Handler {
	return &addStationDataServer{
		svc:       svc,
		transport: newAddStationDataTransport(be),
	}
}

type addPredictDataServer struct {
	svc       service.Storage
	transport *addPredictDataTransport
}

func (s *addPredictDataServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	predictData, err := s.transport.DecodeRequest(r)

	if err == nil {
		err = s.svc.AddPredictedData(r.Context(), predictData)
	}
	s.transport.EncodeResponse(w, err)
}

func addPredictDataHandler(svc service.Storage, be bodyEncodeFunc) http.Handler {
	return &addPredictDataServer{
		svc:       svc,
		transport: newAddPredictDataTransport(be),
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

func getEcoDataListHandler(svc service.Storage, be bodyEncodeFunc) http.Handler {
	return &getEcoDataListServer{
		svc:       svc,
		transport: newGetEcoDataListTransport(be),
	}
}

type getProfilerDataListServer struct {
	svc       service.Storage
	transport *getProfilerDataListTransport
}

func (s *getProfilerDataListServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filter, err := s.transport.DecodeRequest(r)

	var data []service.ProfilerData
	if err == nil {
		data, err = s.svc.GetProfilerDataList(r.Context(), filter)
	}

	s.transport.EncodeResponse(w, data, err)
}

func getProfilerDataListHandler(svc service.Storage, be bodyEncodeFunc) http.Handler {
	return &getProfilerDataListServer{
		svc:       svc,
		transport: newGetProfilerDataListTransport(be),
	}
}

// type predictServer struct {
// 	svc       service.Storage
// 	transport *predictTransport
// }

// func (s *predictServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	filter, err := s.transport.DecodeRequest(r)

// 	if err == nil {
// 		// err = s.svc.Predict(r.Context(), filter)
// 	}

// 	s.transport.EncodeResponse(w, err)
// }

// func predictHandler(svc service.Storage, be bodyEncodeFunc) http.Handler {
// 	return &predictServer{
// 		svc:       svc,
// 		transport: newPredictTransport(be),
// 	}
// }
