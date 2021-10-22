package body

import (
	"encoding/json"
)

// Encode body response by payload and error.
func Encode(payload interface{}, err error) ([]byte, error) {
	body := body{
		IsOk:    err == nil,
		Payload: payload,
	}
	if !body.IsOk {
		body.Payload = err.Error()
	}
	return json.Marshal(body)
}


