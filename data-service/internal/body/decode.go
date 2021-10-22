package body

import (
	"encoding/json"
	"errors"
)

// Decode data to payload.
func Decode(data []byte, payload interface{}) (err error) {
	var body body
	if err = json.Unmarshal(data, &body); err != nil {
		return
	}

	if !body.IsOk {
		err = errors.New(body.Payload.(string))
		return
	}
	data, _ = json.Marshal(body.Payload)

	return json.Unmarshal(data, payload)
}
