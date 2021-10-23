package http

type addStationRequest struct {
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

type stationListResponse struct {
	Stations []stationResponse `json:"station"`
}

type stationResponse struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

type predictDataRequest struct {
	Data []predictData `json:"data"`
}

type predictData struct {
	StationID   string             `json:"station_id"`
	Timestamp   int64              `json:"timestamp"`
	Measurement map[string]float64 `json:"measurement"`
}

type ecoDataListResponse struct {
	Data []ecoDataResponse `json:"data"`
}

type ecoDataResponse struct {
	StationID            string             `json:"station_id"`
	Timestamp            int64              `json:"timestamp"`
	Measurement          map[string]float64 `json:"measurement,omitempty"`
	PredictedMeasurement map[string]float64 `json:"predicted_measurement,omitempty"`
}

type profilerDataListResponse struct {
	Data []profilerDataResponse `json:"data"`
}

type profilerDataResponse struct {
	StationID          string             `json:"station_id"`
	Timestamp          int64              `json:"timestamp"`
	Temperature        map[string]float64 `json:",inline"`
	OutsideTemperature *float64           `json:"outside_temperature,omitempty"`
	WindDirection      *int               `json:"wind_direction,omitempty"`
	WindSpeed          *float64           `json:"wind_speed,omitempty"`
}
