package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/carbocation/interpose"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/topher200/forty-thieves/dal"
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

func (testSuite *MainTestSuite) TestLoginGet() {
	client := http.Client{}
	resp, err := client.Get(testSuite.server.URL + "/login")
	defer resp.Body.Close()
	assert.Nil(testSuite.T(), err)
	assert.Equal(testSuite.T(), 200, resp.StatusCode)
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

func TestLogoutGet(t *testing.T) {
	r, err := http.NewRequest("GET", "/logout", nil)
	assert.Nil(t, err)
	w := runRequest(t, r)
	assert.Equal(t, 302, w.Code)
}

func TestLoginPost(t *testing.T) {
	form := url.Values{
		"Email":    {testUserEmail},
		"Password": {testUserPassword},
	}
	r, err := http.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	assert.Nil(t, err)
	w := runRequest(t, r)
	assert.Equal(t, 302, w.Code)
}

type MainTestSuite struct {
	suite.Suite
	server *httptest.Server
}

func (testSuite *MainTestSuite) SetupSuite() {
	middle := newMiddlewareForTesting(testSuite.T())
	testSuite.server = httptest.NewServer(middle)
}

func deleteTestUser(t *testing.T) {
	userDB := dal.NewUserForTest(t)
	userRow, err := userDB.GetByEmail(nil, testUserEmail)
	assert.Nil(t, err)
	_, err = userDB.DeleteById(nil, userRow.ID)
	assert.Nil(t, err)
}

func (testSuite *MainTestSuite) TearDownSuite() {
	testSuite.server.Close()
	deleteTestUser(testSuite.T())
}

func TestMainSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}
