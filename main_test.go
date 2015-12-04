package main

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
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

type MainTestSuite struct {
	suite.Suite
	server *httptest.Server
	client *http.Client
}

// checkResponse asserts that we didn't err and that our response looks good
func checkResponse(t *testing.T, resp *http.Response, err error) {
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// makeGetRequest makes the request and error-checks the response
func (testSuite *MainTestSuite) makeGetRequest(route string) {
	resp, err := testSuite.client.Get(testSuite.server.URL + route)
	defer resp.Body.Close()
	checkResponse(testSuite.T(), resp, err)
}

func (testSuite *MainTestSuite) TestRootGetBeforeLogin() {
	testSuite.makeGetRequest("/")
}

func (testSuite *MainTestSuite) TestSignupGet() {
	testSuite.makeGetRequest("/signup")
}

func (testSuite *MainTestSuite) TestSignupPost() {
	form := url.Values{
		"Email":         {testUserEmail},
		"Password":      {testUserPassword},
		"PasswordAgain": {testUserPassword},
	}
	resp, err := testSuite.client.PostForm(testSuite.server.URL+"/signup", form)
	checkResponse(testSuite.T(), resp, err)
}

func (testSuite *MainTestSuite) TestRootGetAfterLogin() {
	testSuite.makeGetRequest("/")
}

func (testSuite *MainTestSuite) TestLogoutGet() {
	testSuite.makeGetRequest("/logout")
}

func (testSuite *MainTestSuite) TestLoginGet() {
	testSuite.makeGetRequest("/login")
}

func (testSuite *MainTestSuite) TestLoginPost() {
	form := url.Values{
		"Email":    {testUserEmail},
		"Password": {testUserPassword},
	}
	resp, err := testSuite.client.PostForm(testSuite.server.URL+"/login", form)
	checkResponse(testSuite.T(), resp, err)
}

func (testSuite *MainTestSuite) TestStateGet() {
	resp, err := testSuite.client.Get(testSuite.server.URL + "/state")
	defer resp.Body.Close()
	assert.Nil(testSuite.T(), err)
}

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

func (testSuite *MainTestSuite) SetupSuite() {
	middle := newMiddlewareForTesting(testSuite.T())
	testSuite.server = httptest.NewServer(middle)
	jar, err := cookiejar.New(nil)
	assert.Nil(testSuite.T(), err)
	testSuite.client = &http.Client{Jar: jar}
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
