package predict

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/nECOnetic/data-service/internal/service"
	"golang.org/x/net/http2"
)

type bodyDecodeFunc func(bodyData []byte, payload interface{}) error

type storage interface {
	StoreEcoData(ctx context.Context, dataList []service.EcoData) error
}

// Client for work with predicted module.
type Client struct {
	client *http.Client

	ctx              context.Context
	storage          storage
	decode           bodyDecodeFunc
	predictClientURL string

	logger log.Logger
}

// NewClient ...
func NewClient(
	ctx context.Context,

	storage storage,
	decoode bodyDecodeFunc,
	predictClientURL string,
	logger log.Logger,
) *Client {
	client := &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}
	return &Client{
		ctx:              ctx,
		client:           client,
		storage:          storage,
		decode:           decoode,
		predictClientURL: predictClientURL,
		logger:           logger,
	}
}

// Predict...
// TODO:
// use kafka
func (c *Client) Predict(ctx context.Context, ecoData []service.EcoData, profilerData []service.ProfilerData) error {
	if len(profilerData) == 0 {
		return errProfilerDataIsNotExist
	}
	if len(ecoData) == 0 {
		return errEcoDataIsNotExist
	}

	go func() {
		logger := log.With(c.logger, "method", "Predict", "station", ecoData[0].StationID)
		req := request{
			Data: make([]measurement, 0, len(ecoData)),
		}
		for i, j := 0, 0; i < len(ecoData) && j < len(profilerData); i++ {
			req.Data[i] = measurement{
				Timestamp:   ecoData[i].Timestamp,
				Measurement: ecoData[i].Measurement,
			}
			for ; j < len(profilerData); j++ {
				if ecoData[i].Timestamp == profilerData[j].Timestamp {
					req.Data[i].Temperature = profilerData[j].Temperature
					req.Data[i].OutsideTemperature = profilerData[j].OutsideTemperature
					req.Data[i].WindDirection = profilerData[j].WindDirection
					req.Data[i].WindSpeed = profilerData[j].WindSpeed
				}
			}
		}

		start := time.Now()
		level.Error(logger).Log("msg", "send request to predict service", "data length", len(req.Data))
		res, err := c.sendRequest(c.ctx, req)
		if err != nil {
			level.Error(logger).Log("msg", "send request to predict service", "err", err)
			return
		}
		level.Error(logger).Log("msg", "send request to predict service", "time", time.Since(start).Seconds())

		predictedData := make([]service.EcoData, 0, len(res.Data))
		for _, d := range res.Data {
			predictedData = append(
				predictedData,
				service.EcoData{
					StationID: ecoData[0].StationID,
					Timestamp: d.Timestamp,

					PredictedMeasurement: d.Measurement,
				},
			)
		}
		if err = c.storage.StoreEcoData(c.ctx, predictedData); err != nil {
			level.Error(logger).Log("msg", "store predict data", "err", err)
			return
		}
	}()
	return nil
}

func (c *Client) sendRequest(ctx context.Context, r request) (response, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return response{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.predictClientURL, bytes.NewReader(data))
	if err != nil {
		return response{}, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return response{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response{}, err
	}

	var rr response
	err = c.decode(body, &res)
	return rr, err
}
