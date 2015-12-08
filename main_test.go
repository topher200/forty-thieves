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

// TestUserStory simulates a user performing the following actions:
//  - gets the root page before logging in
//  - gets the signup page
//  - posts to the signup page
//  - gets the logout page
//  - gets the login page
//  - posts to the login page
//  - gets a json /state message
//
// We do this in one function (as opposed to separate Test* functions) since
// each of the tests requires a certain state of the user cookie (user
// exists/doesn't exist, user logged in/out).
func (testSuite *MainTestSuite) TestUserStory() {
	testSuite.makeGetRequest("/")
	testSuite.makeGetRequest("/signup")
	testSuite.signupPost()
	testSuite.makeGetRequest("/")
	testSuite.makeGetRequest("/logout")
	testSuite.makeGetRequest("/login")
	testSuite.loginPost()
	testSuite.stateGet()
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

// signupPost assumes you're signed out and that the user doesn't exist
func (testSuite *MainTestSuite) signupPost() {
	form := url.Values{
		"Email":         {testUserEmail},
		"Password":      {testUserPassword},
		"PasswordAgain": {testUserPassword},
	}
	resp, err := testSuite.client.PostForm(testSuite.server.URL+"/signup", form)
	checkResponse(testSuite.T(), resp, err)
}

// loginPost assumes you're currently signed out and the user exists
func (testSuite *MainTestSuite) loginPost() {
	form := url.Values{
		"Email":    {testUserEmail},
		"Password": {testUserPassword},
	}
	resp, err := testSuite.client.PostForm(testSuite.server.URL+"/login", form)
	checkResponse(testSuite.T(), resp, err)
}

// stateGet assumes you're already signed in
func (testSuite *MainTestSuite) stateGet() {
	testSuite.makeGetRequest("/state")
	// TODO: check that our json looks good
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
