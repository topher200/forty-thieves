package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShowHttp(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	showHttp(w, req)
	assert.Equal(t, 200, w.Code)
	assert.NotEqual(t, "", w.Body.String())
}
