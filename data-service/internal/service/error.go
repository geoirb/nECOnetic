package service

import (
	"errors"
)

var (
	errUnknownType = errors.New("unknown data type")
	errFormateData = errors.New("error formate of data")
)
