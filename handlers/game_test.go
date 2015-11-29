package handlers

import "testing"

func TestHandleStateRequest(t *testing.T) {
	// TODO Set context from base so this doesn't crash
	// req, err := http.NewRequest("GET", "/state", nil)
	// assert.Nil(t, err)
	// w := httptest.NewRecorder()
	// HandleStateRequest(w, req)
	// assert.Equal(t, 200, w.Code)
	// assert.NotEqual(t, "", w.Body.String())
}

func TestHandleMoveRequest(t *testing.T) {
	// jsonStr := []byte(
	// 	`{"FromLocation": "tableau", "FromIndex": 0, "ToLocation": "tableau", "ToIndex": 1 }`)
	// req, err := http.NewRequest("POST", "", bytes.NewBuffer(jsonStr))
	// assert.Nil(t, err)
	// w := httptest.NewRecorder()
	// HandleMoveRequest(w, req)
	// assert.Equal(t, 200, w.Code)
}
