package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleResources(t *testing.T) {
	req, err := http.NewRequest("GET", "/res/cards-png/spades-A.png", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	handleResources(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, w.Header().Get("Content-Type"), "image/png")
}

func TestHandleJavascript(t *testing.T) {
	req, err := http.NewRequest("GET", "/bower_components/fallback/fallback.js", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	handleJavascript(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, w.Header().Get("Content-Type"), "application/javascript")
}

func TestHandleStateRequest(t *testing.T) {
	gameState = NewGame()
	req, err := http.NewRequest("GET", "/state", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	handleStateRequest(w, req)
	assert.Equal(t, 200, w.Code)
	assert.NotEqual(t, "", w.Body.String())
}

func TestHandleMoveRequest(t *testing.T) {
	gameState = NewGame()
	jsonStr := []byte(
		`{"FromLocation": "tableau", "FromIndex": 0, "ToLocation": "tableau", "ToIndex": 1 }`)
	req, err := http.NewRequest("POST", "", bytes.NewBuffer(jsonStr))
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	handleMoveRequest(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestShowHttp(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	showHttp(w, req)
	assert.Equal(t, 200, w.Code)
	assert.NotEqual(t, "", w.Body.String())
}
