package libsolver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/topher200/deck"
	"github.com/topher200/forty-thieves/libgame"
)

func createTestingGameState() libgame.GameState {
	state := libgame.NewGame()
	state.Stock.Cards = append(state.Waste.Cards, deck.Card{Face: "2", Suit: "clubs"})
	state.Stock.Cards = append(state.Foundations[0].Cards, deck.Card{Face: "A", Suit: "clubs"})
	return state
}

func TestNumPiles(t *testing.T) {
	assert.Equal(t, len(allPiles()), 20)
	state := createTestingGameState()
	getPossibleMoves(&state)
	assert.Nil(t, FoundationAvailableCard(&state))
}
