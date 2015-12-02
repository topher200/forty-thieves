package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/carbocation/interpose"
	"github.com/stretchr/testify/assert"
)

const (
	testUserEmail    = "fake@asdf.com"
	testUserPassword = "Password1"
)

func newApplicationForTesting(t *testing.T) *Application {
	app, err := NewApplication(true)
	assert.Nil(t, err)
	return app
}

func newMiddlewareForTesting(t *testing.T) *interpose.Middleware {
	app := newApplicationForTesting(t)
	middle, err := app.middlewareStruct()
	assert.Nil(t, err)
	return middle
}

func runRequest(t *testing.T, r *http.Request) *httptest.ResponseRecorder {
	middle := newMiddlewareForTesting(t)
	w := httptest.NewRecorder()
	middle.ServeHTTP(w, r)
	return w
}

func TestLoginGet(t *testing.T) {
	r, err := http.NewRequest("GET", "/login", nil)
	assert.Nil(t, err)
	w := runRequest(t, r)
	assert.Equal(t, 200, w.Code)
	assert.NotEqual(t, "", w.Body.String())
}

func TestSignupPost(t *testing.T) {
	form := url.Values{
		"Email":         {testUserEmail},
		"Password":      {testUserPassword},
		"PasswordAgain": {testUserPassword},
	}
	r, err := http.NewRequest("POST", "/signup", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	assert.Nil(t, err)
	w := runRequest(t, r)
	assert.Equal(t, 302, w.Code)
}
