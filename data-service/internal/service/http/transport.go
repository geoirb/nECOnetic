package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/nECOnetic/data-service/internal/service"
)

type addStationTransport struct {
	buildResponse buildResponseFunc
}

func newAddStationTransport(
	build buildResponseFunc,
) *addStationTransport {
	return &addStationTransport{
		buildResponse: build,
	}
}

func (*addStationTransport) DecodeRequest(r *http.Request) (s service.Station, err error) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	var req addStationRequest
	if err = json.Unmarshal(body, &req); err != nil {
		return
	}

	s = service.Station{
		Name: req.Name,
		Lat:  req.Lat,
		Lon:  req.Lon,
	}
	return
}

func (t *addStationTransport) EncodeResponse(w http.ResponseWriter, s service.Station, err error) {
	res := stationResponse(s)
	body, _ := t.buildResponse(res, err)
	w.Write(body)
}

type getStationListTransport struct {
	buildResponse buildResponseFunc
}

func newGetStationListTransport(
	build buildResponseFunc,
) *getStationListTransport {
	return &getStationListTransport{
		buildResponse: build,
	}
}

func (t *getStationListTransport) EncodeResponse(w http.ResponseWriter, sList []service.Station, err error) {
	res := stationListResponse{
		Stations: make([]stationResponse, 0, len(sList)),
	}

	for _, s := range sList {
		res.Stations = append(res.Stations, stationResponse(s))
	}

	body, _ := t.buildResponse(res, err)
	w.Write(body)
}

type addStationDataTransport struct {
	buildResponse buildResponseFunc
}

func newAddStationDataTransport(
	build buildResponseFunc,
) *addStationDataTransport {
	return &addStationDataTransport{
		buildResponse: build,
	}
}

func (*addStationDataTransport) DecodeRequest(r *http.Request) (data service.StationData, err error) {
	if err = r.ParseMultipartForm(32 << 20); err != nil {
		return
	}
	v := r.MultipartForm.Value

	station := v["station"]
	if len(station) != 1 {
		err = fmt.Errorf("wrong numbers of station values need: 1 have: %d", len(station))
		return
	}
	data.StationName = station[0]

	dataType := r.Form["type"]
	if len(dataType) != 1 {
		err = fmt.Errorf("wrong numbers of data type values need: 1 have: %d", len(dataType))
		return
	}
	data.Type = dataType[0]

	fileHeader := r.MultipartForm.File["data"]
	if len(fileHeader) != 1 {
		err = fmt.Errorf("wrong numbers of data file values need: 1 have: %d", len(fileHeader))
		return
	}

	data.FileName = fileHeader[0].Filename
	if data.File, err = fileHeader[0].Open(); err != nil {
		return
	}

	return
}

func (t *addStationDataTransport) EncodeResponse(w http.ResponseWriter, err error) {
	body, _ := t.buildResponse(nil, err)
	w.Write(body)
}

type getEcoDataListTransport struct {
	buildResponse buildResponseFunc
}

func newGetEcoDataListTransport(
	build buildResponseFunc,
) *getEcoDataListTransport {
	return &getEcoDataListTransport{
		buildResponse: build,
	}
}

func (*getEcoDataListTransport) DecodeRequest(r *http.Request) (f service.GetEcoData, err error) {
	query := r.URL.Query()

	if stationName, isExist := query["station"]; isExist {
		f.StationName = &stationName[0]
	}

	if timestampFrom, isExist := query["timestamp_from"]; isExist {
		var tsFrom int64
		if tsFrom, err = strconv.ParseInt(timestampFrom[0], 10, 64); err != nil {
			err = fmt.Errorf("parse timestamp_from: %s", err)
			return
		}
		f.TimestampFrom = &tsFrom
	}

	if timestampTill, isExist := query["timestamp_till"]; isExist {
		var tsTil int64
		if tsTil, err = strconv.ParseInt(timestampTill[0], 10, 64); err != nil {
			err = fmt.Errorf("parse timestamp_till: %s", err)
			return
		}
		f.TimestampTill = &tsTil
	}

	f.Measurements = query["measurement"]
	return
}

func (t *getEcoDataListTransport) EncodeResponse(w http.ResponseWriter, data []service.EcoData, err error) {
	res := ecoDataListResponse{
		Data: make([]ecoDataResponse, 0, len(data)),
	}

	for _, d := range data {
		res.Data = append(res.Data, ecoDataResponse(d))
	}

	body, _ := t.buildResponse(res, err)
	w.Write(body)
}
