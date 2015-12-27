package libgame

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/topher200/deck"
)

func TestNewGame(t *testing.T) {
	state := NewGame()
	assert.NotEmpty(t, state.Stock.Cards)
	for _, foundation := range state.Foundations {
		assert.Empty(t, foundation.Cards)
	}
	for _, tableau := range state.Tableaus {
		assert.Len(t, tableau.Cards, 4)
	}
}

func TestPopFromStock(t *testing.T) {
	state := NewGame()
	numCards := len(state.Stock.Cards)
	card, err := state.popFromStock()
	assert.Nil(t, err)
	assert.NotEqual(t, card, deck.Card{})
	assert.Equal(t, numCards-1, len(state.Stock.Cards))
}

func TestMoveCard(t *testing.T) {
	// Make sure that the size of each deck has inc/decreased, and that the moved
	// card is now at the bottom of deck #2
	state := NewGame()
	deckFrom := state.Tableaus[0]
	deckFromLen := len(deckFrom.Cards)
	cardToMove := deckFrom.Cards[deckFromLen-1]
	deckTo := state.Tableaus[1]
	deckToLen := len(deckTo.Cards)

	err := state.MoveCard(&deckFrom, &deckTo)

	assert.Nil(t, err)
	assert.Len(t, deckFrom.Cards, deckFromLen-1)
	assert.Len(t, deckTo.Cards, deckToLen+1)
	cardInNewSpot := deckTo.Cards[len(deckTo.Cards)-1]
	assert.Equal(t, cardToMove, cardInNewSpot)
}

func TestFlipStock(t *testing.T) {
	state := NewGame()
	stockLenStart := len(state.Stock.Cards)
	wasteLenStart := len(state.Waste.Cards)

	err := state.FlipStock()
	assert.Nil(t, err)
	assert.Len(t, state.Stock.Cards, stockLenStart-1)
	assert.Len(t, state.Waste.Cards, wasteLenStart+1)
}

func TestScore(t *testing.T) {
	state := NewGame()
	assert.Equal(t, 104, state.Score)
}
