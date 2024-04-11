//go:build test_all || test_e2e

package main

import (

	// New import
	"net/http"
	"testing"

	"github.com/mixnblend/snippetbox/internal/assert"
)

func TestPingE2E(t *testing.T) {
	// Given ... we have an application with a structured logger which discards everthing.
	app := newTestApplication(t)

	// And ... we have created a new test server
	testServer := newTestServer(t, app.routes())
	defer testServer.Close()

	// When ...we call the ping endpoint
	code, _, body := testServer.get(t, "/ping")

	// Then ... the OK response should be returned as expected
	assert.Equal(t, code, http.StatusOK)
	// And ... a body with the payload "OK" should be returned as expected
	assert.Equal(t, body, "OK")
}

func TestSnippetViewE2E(t *testing.T) {
	// Given ... we have an application with a structured logger which discards everthing.
	app := newTestApplication(t)

	// And ... we have created a new test server
	testServer := newTestServer(t, app.routes())
	defer testServer.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond...",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/snippet/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tableTest := range tests {
		t.Run(tableTest.name, func(t *testing.T) {
			// When ... we call our path
			code, _, body := testServer.get(t, tableTest.urlPath)

			// Then ... the HTTP status code should be returned as expected
			assert.Equal(t, code, tableTest.wantCode)

			// And ... the body of the response should be returned as expected.
			if tableTest.wantBody != "" {
				assert.StringContains(t, body, tableTest.wantBody)
			}
		})
	}
}
