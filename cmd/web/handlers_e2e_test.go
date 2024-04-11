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
