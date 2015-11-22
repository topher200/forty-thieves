package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
