package main

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Sirupsen/logrus"
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
	testSuite.newgamePost()
	testSuite.stateGet()
	testSuite.flipStockPost()
	testSuite.movePost()
	// TODO: compare this gamestate to the gamestate before the last move?
	testSuite.undoMovePost()
}

// checkResponse asserts that we didn't err and that our response looks good
func checkResponse(t *testing.T, resp *http.Response, err error) {
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// makeGetRequest makes the request and error-checks the response
func (testSuite *MainTestSuite) makeGetRequest(route string) []byte {
	resp, err := testSuite.client.Get(testSuite.server.URL + route)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(testSuite.T(), err)
	defer resp.Body.Close()
	checkResponse(testSuite.T(), resp, err)
	return body
}

// signupPost assumes you're signed out and that the user doesn't exist
func (testSuite *MainTestSuite) signupPost() {
	form := url.Values{
		"Email":         {testUserEmail},
		"Password":      {testUserPassword},
		"PasswordAgain": {testUserPassword},
	}
	resp, err := testSuite.client.PostForm(testSuite.server.URL+"/signup", form)
	defer resp.Body.Close()
	checkResponse(testSuite.T(), resp, err)
}

// loginPost assumes you're currently signed out and the user exists
func (testSuite *MainTestSuite) loginPost() {
	form := url.Values{
		"Email":    {testUserEmail},
		"Password": {testUserPassword},
	}
	resp, err := testSuite.client.PostForm(testSuite.server.URL+"/login", form)
	defer resp.Body.Close()
	checkResponse(testSuite.T(), resp, err)
}

// newgamePost assumes that you're signed in
func (testSuite *MainTestSuite) newgamePost() {
	resp, err := testSuite.client.Post(testSuite.server.URL+"/newgame", "text/json", nil)
	defer resp.Body.Close()
	checkResponse(testSuite.T(), resp, err)
}

// flipStockPost assumes that you're signed in
func (testSuite *MainTestSuite) flipStockPost() {
	resp, err := testSuite.client.Post(
		testSuite.server.URL+"/flipstock", "text/json", nil)
	defer resp.Body.Close()
	checkResponse(testSuite.T(), resp, err)
}

// movePost assumes that you're signed in
func (testSuite *MainTestSuite) movePost() {
	form := url.Values{
		"FromLocation": {"tableau"},
		"FromIndex":    {"0"},
		"ToLocation":   {"tableau"},
		"ToIndex":      {"1"},
	}
	resp, err := testSuite.client.PostForm(testSuite.server.URL+"/move", form)
	defer resp.Body.Close()
	checkResponse(testSuite.T(), resp, err)
}

// undoMovePost assumes that you're signed in
func (testSuite *MainTestSuite) undoMovePost() {
	resp, err := testSuite.client.Post(testSuite.server.URL+"/undomove", "text/json", nil)
	defer resp.Body.Close()
	checkResponse(testSuite.T(), resp, err)
}

// stateGet assumes you're already signed in
func (testSuite *MainTestSuite) stateGet() {
	bodyText := string(testSuite.makeGetRequest("/state"))
	// Check that our response contains one of the card pile names we expect
	assert.True(testSuite.T(), strings.Contains(bodyText, "Stock"))
}

func newApplicationForTesting(t *testing.T) *Application {
	app, err := NewApplication(true)
	assert.Nil(t, err)
	return app
}

func newMiddlewareForTesting(t *testing.T) *interpose.Middleware {
	app := newApplicationForTesting(t)
	logWriter := logrus.New().Writer()
	defer logWriter.Close()
	middle, err := app.middlewareStruct(logWriter)
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
