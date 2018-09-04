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
)

type MainTestSuite struct {
	suite.Suite
	server *httptest.Server
	client *http.Client
}

// TestUserStory simulates a user performing the following actions:
//  - gets the root page
//  - posts to create a new game
//  - gets a json /state message
//  - posts to flip the stock
//  - gets a json /state message
//
// TODO: We do this in one function (as opposed to separate Test* functions)
// since some tests require setup (like a game to be created).
func (testSuite *MainTestSuite) TestUserStory() {
	testSuite.makeGetRequest("/")
	gameStateID := testSuite.newgamePost()
	testSuite.stateGet(gameStateID)
	testSuite.flipStockPost(gameStateID)
	testSuite.stateGet(gameStateID)
	testSuite.movePost(gameStateID)
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

// newgamePost confirms that creating a new game works successfully
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

// addGameStateIdToURL is a helper function for structuring our request URLs
func addGameStateIdToURL(url string, gameStateID uuid.UUID) string {
	return fmt.Sprintf("%s?gameStateID=%s", url, gameStateID.String())
}

// flipStockPost tests that we can flip the stock card
func (testSuite *MainTestSuite) flipStockPost(gameStateID uuid.UUID) {
	resp, err := testSuite.client.Post(
		addGameStateIdToURL(testSuite.server.URL+"/flipstock", gameStateID),
		"text/json", nil)
	defer resp.Body.Close()
	checkResponse(testSuite.T(), resp, err)
}

// movePost tests that we can move a card from one pile to another
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

	// we're not guaranteed to have a move available. we just check that
	// either the request completed or that we threw a validation error
	if resp.StatusCode == 200 {
		checkResponse(testSuite.T(), resp, err)
	} else {
		assert.Equal(testSuite.T(), 400, resp.StatusCode)
	}
}

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

func (testSuite *MainTestSuite) TearDownSuite() {
	testSuite.server.Close()
}

func TestMainSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}
