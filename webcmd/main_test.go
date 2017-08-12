package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/carbocation/interpose"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/topher200/forty-thieves/libdb"
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
//  - posts to create a new game
//  - gets a json /state message
//  - posts to flip the stock
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

	gameStateID := testSuite.newgamePost()
	testSuite.stateGet(gameStateID)
	testSuite.flipStockPost(gameStateID)
	testSuite.stateGet(gameStateID)

	// TODO(topher): test move post by checking that we don't get a 5xx error
	// testSuite.movePost()
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
func (testSuite *MainTestSuite) newgamePost() uuid.UUID {
	resp, err := testSuite.client.Post(testSuite.server.URL+"/newgame", "text/json", nil)
	defer resp.Body.Close()
	checkResponse(testSuite.T(), resp, err)

	// pull out the gameStateID
	type Response struct {
		GameStateID uuid.UUID
	}
	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.Nil(testSuite.T(), err)
	return response.GameStateID
}

func addGameStateIdToURL(url string, gameStateID uuid.UUID) string {
	return fmt.Sprintf("%s?gameStateID=%s", url, gameStateID.String())
}

// flipStockPost assumes that you're signed in
func (testSuite *MainTestSuite) flipStockPost(gameStateID uuid.UUID) {
	resp, err := testSuite.client.Post(
		addGameStateIdToURL(testSuite.server.URL+"/flipstock", gameStateID),
		"text/json", nil)
	defer resp.Body.Close()
	checkResponse(testSuite.T(), resp, err)
}

// movePost assumes that you're signed in
func (testSuite *MainTestSuite) movePost(gameStateID uuid.UUID) {
	form := url.Values{
		"FromPile":  {"tableau"},
		"FromIndex": {"0"},
		"ToPile":    {"tableau"},
		"ToIndex":   {"1"},
	}
	resp, err := testSuite.client.PostForm(
		addGameStateIdToURL(testSuite.server.URL+"/move", gameStateID), form)
	defer resp.Body.Close()
	checkResponse(testSuite.T(), resp, err)
}

// stateGet assumes you're already signed in
func (testSuite *MainTestSuite) stateGet(gameStateID uuid.UUID) {
	bodyText := string(testSuite.makeGetRequest(
		addGameStateIdToURL("/state", gameStateID)))
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
	userDB := libdb.NewUserDBForTest(t)
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
