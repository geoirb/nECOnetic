package predicted

import (
	"github.com/nECOnetic/data-service/internal/service"
)


type storage interface {
	StoreEcoData(ctx context.Context, dataList []EcoData) error
}

type Client struct {
	predictedCh chan []service.EcoData

	measurementCache map[string][]service.EcoData
	profilerCache    []service.ProfilerData
}

func NewClient() *Client{
	return 
}

func (c *Client) 
