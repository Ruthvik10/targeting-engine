package jsonutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteJSON(t *testing.T) {
	// Define a struct for testing JSON output
	type testResponse struct {
		Message string `json:"message"`
	}

	t.Run("should write valid JSON response", func(t *testing.T) {
		// Create a response recorder
		rec := httptest.NewRecorder()

		// Define the response to send
		response := testResponse{Message: "Success"}

		// Call WriteJSON with the response
		WriteJSON(rec, http.StatusOK, response)

		// Check the status code
		assert.Equal(t, http.StatusOK, rec.Code)

		// Check the content-type header
		assert.Equal(t, "application/json", rec.Header().Get("content-type"))

		// Parse the JSON response
		var result testResponse
		err := json.NewDecoder(rec.Body).Decode(&result)
		assert.NoError(t, err)

		// Check the response body
		assert.Equal(t, response.Message, result.Message)
	})

	t.Run("should return 500 with invalid JSON", func(t *testing.T) {
		// Create a response recorder
		rec := httptest.NewRecorder()

		// Pass an invalid data type (channel) to WriteJSON
		type InvalidType struct {
			Ch chan int `json:"ch"`
		}

		WriteJSON(rec, http.StatusInternalServerError, InvalidType{Ch: make(chan int)})

		// Ensure the response is 500 as JSON encoding should fail
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestWriteError(t *testing.T) {
	t.Run("should write error message in JSON format", func(t *testing.T) {
		// Create a response recorder
		rec := httptest.NewRecorder()

		// Define the error message
		errorMessage := "something went wrong"

		// Call WriteError
		WriteError(rec, http.StatusBadRequest, errorMessage)

		// Check the status code
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// Check the content-type header
		assert.Equal(t, "application/json", rec.Header().Get("content-type"))

		// Parse the JSON response
		var result map[string]string
		err := json.NewDecoder(rec.Body).Decode(&result)
		assert.NoError(t, err)

		// Check the error message
		assert.Equal(t, errorMessage, result["error"])
	})
}
