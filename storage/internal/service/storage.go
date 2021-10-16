package service

import (
	"context"
)

type Storage interface{
	AddDataFromStation(ctx context.Context, in StationData) error
}
