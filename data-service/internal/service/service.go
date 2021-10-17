package service

import (
	"bufio"
	"context"
	"errors"
	"io"
	"regexp"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

type storage interface {
	StoreStation(ctx context.Context, st Station) (Station, error)
	StoreProfilerData(ctx context.Context, dataList []ProfilerData) error
	StoreEcoData(ctx context.Context, dataList []EcoData) error

	LoadStation(ctx context.Context, f StationFilter) ([]Station, error)
}

// TODO:
// Logging

type service struct {
	storage storage

	dataHandler map[string]func(context.Context, string, io.Reader) (err error)
}

// New returns Storage interface for work with storage.
func New(
	storage storage,
) Storage {
	s := &service{
		storage: storage,
	}

	s.dataHandler = map[string]func(context.Context, string, io.Reader) (err error){
		"eco":         s.ecoDataHandler,
		"wind":        s.windHandler,
		"temperature": s.temperatureHandler,
	}

	return s
}

// AddDataFromStation to storage.
func (s *service) AddDataFromStation(ctx context.Context, in StationData) error {
	stationFilter := StationFilter{
		Name: in.StationName,
	}

	stations, err := s.storage.LoadStation(ctx, stationFilter)
	if err != nil {
		return err
	}

	h, isExist := s.dataHandler[in.Type]
	if !isExist {
		return errors.New("unknown data type")
	}

	return h(ctx, stations[0].ID, in.File)
}

func (s *service) ecoDataHandler(ctx context.Context, stationID string, r io.Reader) (err error) {
	in, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}

	name := in.GetSheetName(0)
	rows, err := in.GetRows(name)
	if err != nil {
		return err
	}

	hRow := rows[0]
	dataRows := rows[1:]

	dataList := make([]EcoData, 0, len(dataRows))

	loc := time.Now().Location()
	for _, d := range dataRows {
		el := EcoData{
			StationID:   stationID,
			Measurement: make(map[string]float64),
		}

		for i := 1; i < len(hRow) && i < len(d); i++ {
			if len(d[i]) != 0 && len(hRow[i]) != 0 {
				if el.Measurement[hRow[i]], err = strconv.ParseFloat(d[i], 64); err != nil {
					return
				}
			}
		}
		if len(el.Measurement) == 0 || len(d[0]) == 0 {
			continue
		}

		dt, err := time.Parse("02/01/2006 15:04", d[0])
		if err != nil {
			return err
		}
		el.Datatime = time.Date(
			dt.Year(),
			dt.Month(),
			dt.Day(),
			dt.Hour(), dt.Minute(),
			0,
			0,
			loc,
		).Unix()

		dataList = append(dataList, el)
	}
	return s.storage.StoreEcoData(ctx, dataList)
}

func (s *service) windHandler(ctx context.Context, stationID string, r io.Reader) (err error) {
	in, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}

	name := in.GetSheetName(0)
	rows, err := in.GetRows(name)
	if err != nil {
		return err
	}

	dataRows := rows[2:]

	dataList := make([]ProfilerData, 0, len(dataRows))

	loc := time.Now().Location()
	for _, d := range dataRows {
		el := ProfilerData{
			StationID: stationID,
		}

		var dt time.Time
		dt, err = time.Parse("02/01/2006 15:04", d[0])
		if err != nil {
			continue
		}
		el.Datatime = time.Date(
			dt.Year(),
			dt.Month(),
			dt.Day(),
			dt.Hour(), dt.Minute(),
			0,
			0,
			loc,
		).Unix()

		var windDirection int
		windDirection, err = strconv.Atoi(d[1])
		if err != nil {
			continue
		}
		if windDirection < 0 || windDirection > 360 {
			continue
		}
		el.WindDirection = &windDirection

		var windSpeed int
		windSpeed, err = strconv.Atoi(d[2])
		if err != nil {
			continue
		}
		el.WindSpeed = &windSpeed

		dataList = append(dataList, el)
	}
	return s.storage.StoreProfilerData(ctx, dataList)
}

func (s *service) temperatureHandler(ctx context.Context, stationID string, r io.Reader) (err error) {
	hRegexp := regexp.MustCompile(`data.*time.*OutsideTemperature.*Quality`)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if !hRegexp.Match(scanner.Bytes()) {
			continue
		}

	}

	return s.storage.StoreProfilerData(ctx, nil)
}
