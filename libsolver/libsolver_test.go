package libsolver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/topher200/deck"
	"github.com/topher200/forty-thieves/libgame"
)

func createTestingGameState() libgame.GameState {
	game := libgame.Game{1}
	state := libgame.DealNewGame(game)
	state.Stock.Cards = append(state.Waste.Cards, deck.Card{Face: "2", Suit: "clubs"})
	state.Stock.Cards = append(state.Foundations[0].Cards, deck.Card{Face: "A", Suit: "clubs"})
	return state
}

func TestNumPiles(t *testing.T) {
	assert.Equal(t, len(allPiles()), 20)
}

func TestGetPossibleMovesReturnsAMove(t *testing.T) {
	state := createTestingGameState()
	moves := GetPossibleMoves(&state)
	assert.NotEmpty(t, moves)
}

func TestFoundationACard(t *testing.T) {
	state := createTestingGameState()
	assert.Nil(t, FoundationAvailableCard(&state))
}
