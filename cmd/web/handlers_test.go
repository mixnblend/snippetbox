package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestPingE2E(t *testing.T) {

	endToEndTest(t)
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
	endToEndTest(t)
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

func TestUserSignupE2E(t *testing.T) {
	endToEndTest(t)

	// Given ... we have an application with a structured logger which discards everthing.
	app := newTestApplication(t)

	// And ... we have created a new test server
	testServer := newTestServer(t, app.routes())
	defer testServer.Close()

	// and ... we have called our user signup route and extracted the csrf token
	_, _, body := testServer.get(t, "/user/signup")
	validCSRFToken := extractCSRFToken(t, body)

	type TestCase struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantFormTag  string
	}

	const (
		validName     = "Bob"
		validPassword = "validPa$$word"
		validEmail    = "bob@example.com"
		formTag       = "<form action='/user/signup' method='POST' novalidate>"
	)

	validSubmission := TestCase{
		name:         "Valid submission",
		userName:     validName,
		userEmail:    validEmail,
		userPassword: validPassword,
		csrfToken:    validCSRFToken,
		wantCode:     http.StatusSeeOther,
	}

	invalidCSRFToken := validSubmission
	invalidCSRFToken.name = "Invalid CSRF Token"
	invalidCSRFToken.csrfToken = "wrongToken"
	invalidCSRFToken.wantCode = http.StatusBadRequest

	emptyName := validSubmission
	emptyName.name = "Empty name"
	emptyName.userName = ""
	emptyName.wantCode = http.StatusUnprocessableEntity

	emptyEmail := validSubmission
	emptyEmail.name = "Empty email"
	emptyEmail.userEmail = ""
	emptyEmail.wantCode = http.StatusUnprocessableEntity

	emptyPassword := validSubmission
	emptyPassword.name = "Empty password"
	emptyPassword.userPassword = ""
	emptyPassword.wantCode = http.StatusUnprocessableEntity

	invalidEmail := validSubmission
	invalidEmail.name = "Invalid email"
	invalidEmail.userEmail = "bob@example."
	invalidEmail.wantCode = http.StatusUnprocessableEntity

	shortPassword := validSubmission
	shortPassword.name = "Short password"
	shortPassword.userEmail = "pa$$"
	shortPassword.wantCode = http.StatusUnprocessableEntity

	duplicateEmail := validSubmission
	duplicateEmail.name = "Duplicate email"
	duplicateEmail.userEmail = "dupe@example.com"
	duplicateEmail.wantCode = http.StatusUnprocessableEntity

	tests := []TestCase{
		validSubmission,
		invalidCSRFToken,
		emptyName,
		emptyEmail,
		emptyPassword,
		invalidEmail,
		shortPassword,
		duplicateEmail,
	}

	for _, tableTest := range tests {
		t.Run(tableTest.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tableTest.userName)
			form.Add("email", tableTest.userEmail)
			form.Add("password", tableTest.userPassword)
			form.Add("csrf_token", tableTest.csrfToken)

			code, _, body := testServer.postForm(t, "/user/signup", form)

			assert.Equal(t, code, tableTest.wantCode)

			if tableTest.wantFormTag != "" {
				assert.StringContains(t, body, tableTest.wantFormTag)
			}
		})
	}
}
