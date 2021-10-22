package predict

type request struct {
	Data []measurement `json:"data"`
}

type measurement struct {
	Timestamp          int64              `json:"timestamp"`
	Measurement        map[string]float64 `json:"measurement"`
	Temperature        map[string]float64 `json:"temperature"`
	OutsideTemperature *float64           `json:"outside_temperature"`
	WindDirection      *int               `json:"wind_direction"`
	WindSpeed          *float64           `json:"wind_speed"`
}

type response struct {
	Data []predicted `json:"data"`
}

type predicted struct {
	Timestamp   int64              `json:"timestamp"`
	Measurement map[string]float64 `json:"measurement"`
}
