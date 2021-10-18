package service

import (
	"bufio"
	"context"
	"errors"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

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

		el.Timestamp, err = parseTime(d[0])
		if err != nil {
			return err
		}

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

	for _, d := range dataRows {
		el := ProfilerData{
			StationID: stationID,
		}

		el.Timestamp, err = parseTime(d[0])
		if err != nil {
			return err
		}

		var windDirection int
		windDirection, err = strconv.Atoi(d[1])
		if err != nil {
			return err
		}

		if windDirection < 0 || windDirection > 360 {
			continue
		}
		el.WindDirection = &windDirection

		var windSpeed int
		windSpeed, err = strconv.Atoi(d[2])
		if err != nil {
			return err
		}

		el.WindSpeed = &windSpeed

		dataList = append(dataList, el)
	}
	return s.storage.StoreProfilerData(ctx, dataList)
}

var (
	headerRegexp = regexp.MustCompile(`data\stime.*OutsideTemperature\sQuality`)
	dataRegexp   = regexp.MustCompile(`^([\d]{2}\/[\d]{2}\/[\d]{4} [012]\d:[0-5]\d):[0-5]\d(\s([-]?[0-9,]+))+$`)

	dateRegexp = regexp.MustCompile(`^([\d]{2}\/[\d]{2}\/[\d]{4} [012]\d:[0-5]\d)`)
)

func (s *service) temperatureHandler(ctx context.Context, stationID string, r io.Reader) error {
	var hights []string

	dataList := make([]ProfilerData, 0, 288)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if headerRegexp.Match(scanner.Bytes()) {
			hights = parseHights(scanner.Text())
		}
		if dataRegexp.Match(scanner.Bytes()) {
			timestemp, err := parseTime(string(dateRegexp.Find(scanner.Bytes())))
			if err != nil {
				return err
			}

			measurements := parseDigits(scanner.Text())
			if len(measurements)-2 != len(hights) {
				// TODO:
				// Log
				return errors.New("error formate of data")
			}
			el := ProfilerData{
				StationID:          stationID,
				Timestamp:          timestemp,
				OutsideTemperature: &measurements[len(measurements)-2],
				Temperature:        make(map[string]float64),
			}
			for i, h := range hights {
				el.Temperature[h] = measurements[i]
			}

			dataList = append(dataList, el)
		}
	}

	return s.storage.StoreProfilerData(ctx, dataList)
}

var digitRegexp = regexp.MustCompile(`\t([-]?[0-9]+,?[0-9]*)`)

func parseDigits(str string) []float64 {
	var d []float64
	if submatch := digitRegexp.FindAllStringSubmatch(str, -1); len(submatch) != 0 {
		d = make([]float64, 0, len(submatch))
		for _, sm := range submatch {
			cd, _ := strconv.ParseFloat(strings.Replace(sm[1], ",", ".", -1), 64)
			d = append(d, cd)
		}
	}
	return d
}

func parseHights(str string) []string {
	var d []string
	if submatch := digitRegexp.FindAllStringSubmatch(str, -1); len(submatch) != 0 {
		d = make([]string, 0, len(submatch))
		for _, sm := range submatch {
			d = append(d, sm[1])
		}
	}
	return d
}

var loc = time.Now().Location()

func parseTime(str string) (int64, error) {
	dt, err := time.Parse("02/01/2006 15:04", str)
	if err != nil {
		return 0, err
	}
	return time.Date(
		dt.Year(),
		dt.Month(),
		dt.Day(),
		dt.Hour(), dt.Minute(),
		0,
		0,
		loc,
	).Unix(), nil
}
