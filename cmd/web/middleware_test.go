package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mixnblend/snippetbox/internal/assert"
)

func TestCommonHeaders(t *testing.T) {
	// given ... we have a response recorder
	responseRecorder := httptest.NewRecorder()

	// and ... we have a dummy http request
	dummyRequest, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// and ... we have created a mock http handler that we can pass to our common headers middleware which writes a 200 reponse
	// and an "OK" response body
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// when ... we call our middleware function passing it the handler
	commonHeaders(next).ServeHTTP(responseRecorder, dummyRequest)

	response := responseRecorder.Result()

	// then ... the Content-Security-Policy header should be set as expected
	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, response.Header.Get("Content-Security-Policy"), expectedValue)

	// ... the Referrer-Policy header should have been set as expected.
	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, response.Header.Get("Referrer-Policy"), expectedValue)

	// ... The X-Content-Type-Options header should have been set as expected.
	expectedValue = "nosniff"
	assert.Equal(t, response.Header.Get("X-Content-Type-Options"), expectedValue)

	// ... The X-Frame-Options header should have been set as expected.
	expectedValue = "deny"
	assert.Equal(t, response.Header.Get("X-Frame-Options"), expectedValue)

	// ... The X-XSS-Protection header should have been set as expected.
	expectedValue = "0"
	assert.Equal(t, response.Header.Get("X-XSS-Protection"), expectedValue)

	// ... The Server header should have been set as expected.
	expectedValue = "Go"
	assert.Equal(t, response.Header.Get("Server"), expectedValue)

	// ... the middleware should correctly call the next handler and the status code and body should be set as expected
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
