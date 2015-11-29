package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/carbocation/interpose"
	"github.com/stretchr/testify/assert"
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

func TestGetLogin(t *testing.T) {
	middle := newMiddlewareForTesting(t)
	r, err := http.NewRequest("GET", "/login", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	middle.ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)
	assert.NotEqual(t, "", w.Body.String())
}
