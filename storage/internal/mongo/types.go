package mongo

import (
	"time"
)

type metrics struct {
	Date      time.Time `bson:"status"`
CO float64 
}
