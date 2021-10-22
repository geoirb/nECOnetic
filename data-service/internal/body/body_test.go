package body

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testData = "test-data"
	errTest  = errors.New("test-error")
)

func TestEncode(t *testing.T) {
	t.Run("success payload", func(t *testing.T) {
		payload := struct {
			Data string `json:"data"`
		}{
			Data: testData,
		}

		response := body{
			IsOk:    true,
			Payload: payload,
		}
		expectedData, err := json.Marshal(&response)
		assert.NotNil(t, expectedData)
		assert.NoError(t, err)

		actualData, err := Encode(payload, nil)
		assert.NotNil(t, actualData)
		assert.NoError(t, err)
		assert.Equal(t, expectedData, actualData)
	})

	t.Run("error", func(t *testing.T) {
		payload := struct {
			Data string `json:"data"`
		}{
			Data: testData,
		}

		response := body{
			IsOk:    false,
			Payload: errTest.Error(),
		}
		expectedData, err := json.Marshal(&response)
		assert.NotNil(t, expectedData)
		assert.NoError(t, err)

		actualData, err := Encode(payload, errTest)
		assert.NotNil(t, actualData)
		assert.NoError(t, err)
		assert.Equal(t, expectedData, actualData)
	})
}

type field struct {
	Filed1 string `json:"field_1"`
	Filed2 string `json:"field_2"`
}
type testPayload struct {
	ID    int     `json:"id"`
	Filed []field `json:"field"`
}

func TestDecode(t *testing.T) {
	expectedPayload := testPayload{
		ID: 1,
		Filed: []field{
			{
				Filed1: "field_1",
				Filed2: "field_2",
			},
			{
				Filed1: "field_3",
				Filed2: "field_4",
			},
		},
	}
	expectedData := body{
		IsOk:    true,
		Payload: expectedPayload,
	}

	body, err := json.Marshal(expectedData)
	assert.NoError(t, err)

	var actualData testPayload
	err = Decode(body, &actualData)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, actualData)
}
