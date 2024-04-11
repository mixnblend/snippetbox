//go:build test_all || test_unit

package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mixnblend/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	// Given ... We have an http response recorder
	responseRecorder := httptest.NewRecorder()

	// And ... a dummy http request
	dummyRequest, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// When ... we call our handler
	ping(responseRecorder, dummyRequest)

	response := responseRecorder.Result()

	// Then ... the status code returned should be 200 as expected
	assert.Equal(t, response.StatusCode, http.StatusOK)

	// And ... the response body returned should be "OK"
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
