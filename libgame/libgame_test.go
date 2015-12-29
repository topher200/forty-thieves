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
	state.Tableaus[0].Cards = []deck.Card{
		deck.Card{Face: deck.KING, Suit: deck.CLUB},
		deck.Card{Face: deck.QUEEN, Suit: deck.CLUB},
		deck.Card{Face: deck.JACK, Suit: deck.CLUB},
		deck.Card{Face: deck.NINE, Suit: deck.CLUB}}
	deckFrom := &state.Tableaus[0]
	deckFromLen := len(deckFrom.Cards)
	cardToMove := deckFrom.Cards[deckFromLen-1]
	state.Tableaus[1].Cards = []deck.Card{
		deck.Card{Face: deck.KING, Suit: deck.CLUB},
		deck.Card{Face: deck.QUEEN, Suit: deck.CLUB},
		deck.Card{Face: deck.JACK, Suit: deck.CLUB},
		deck.Card{Face: deck.TEN, Suit: deck.CLUB}}
	deckTo := &state.Tableaus[1]
	deckToLen := len(deckTo.Cards)

	err := state.MoveCard(MoveRequest{tableau, 0, tableau, 1})

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

func TestIsMoveLegal(t *testing.T) {
	// move to waste and stock
	assert.Error(t, IsMoveLegal(
		tableau,
		deck.Deck{Cards: []deck.Card{deck.Card{Face: deck.KING, Suit: deck.CLUB}}},
		waste,
		deck.Deck{}),
		"moving to waste is illegal")
	assert.Error(t, IsMoveLegal(
		tableau,
		deck.Deck{Cards: []deck.Card{deck.Card{Face: deck.KING, Suit: deck.CLUB}}},
		stock,
		deck.Deck{}),
		"moving to stock is illegal")

	// move to foundation
	assert.Nil(t, IsMoveLegal(
		tableau,
		deck.Deck{Cards: []deck.Card{deck.Card{Face: deck.ACE, Suit: deck.CLUB}}},
		foundation,
		deck.Deck{}),
		"moving to empty foundation with ace is OK")
	assert.Error(t, IsMoveLegal(
		tableau,
		deck.Deck{Cards: []deck.Card{deck.Card{Face: deck.KING, Suit: deck.CLUB}}},
		foundation,
		deck.Deck{}),
		"moving non-ace to empty foundation is illegal")
	assert.Error(t, IsMoveLegal(
		tableau,
		deck.Deck{Cards: []deck.Card{deck.Card{Face: deck.ACE, Suit: deck.CLUB}}},
		foundation,
		deck.Deck{Cards: []deck.Card{deck.Card{Face: deck.KING, Suit: deck.CLUB}}}),
		"moving ace to populated foundation is illegal")
	assert.Error(t, IsMoveLegal(
		tableau,
		deck.Deck{Cards: []deck.Card{deck.Card{Face: deck.TEN, Suit: deck.CLUB}}},
		foundation,
		deck.Deck{Cards: []deck.Card{deck.Card{Face: deck.JACK, Suit: deck.CLUB}}}),
		"moving ten onto jack in foundation is illegal")
	assert.Nil(t, IsMoveLegal(
		tableau,
		deck.Deck{Cards: []deck.Card{deck.Card{Face: deck.TWO, Suit: deck.CLUB}}},
		foundation,
		deck.Deck{Cards: []deck.Card{deck.Card{Face: deck.ACE, Suit: deck.CLUB}}}),
		"moving two on top of ace in foundation is OK")
	assert.Nil(t, IsMoveLegal(
		tableau,
		deck.Deck{Cards: []deck.Card{deck.Card{Face: deck.THREE, Suit: deck.CLUB}}},
		foundation,
		deck.Deck{Cards: []deck.Card{
			deck.Card{Face: deck.ACE, Suit: deck.CLUB},
			deck.Card{Face: deck.TWO, Suit: deck.CLUB}}}),
		"moving three on top of ace/two in foundation is OK")

	// moving to tableaus
	assert.Nil(t, IsMoveLegal(
		tableau,
		deck.Deck{Cards: []deck.Card{deck.Card{Face: deck.THREE, Suit: deck.CLUB}}},
		tableau,
		deck.Deck{}),
		"moving to an empty tableau is OK")
	assert.Nil(t, IsMoveLegal(
		tableau,
		deck.Deck{Cards: []deck.Card{deck.Card{Face: deck.TEN, Suit: deck.CLUB}}},
		tableau,
		deck.Deck{Cards: []deck.Card{deck.Card{Face: deck.JACK, Suit: deck.CLUB}}}),
		"moving ten onto jack in tableau is ok")
	assert.Nil(t, IsMoveLegal(
		tableau,
		deck.Deck{Cards: []deck.Card{deck.Card{Face: deck.NINE, Suit: deck.CLUB}}},
		tableau,
		deck.Deck{Cards: []deck.Card{
			deck.Card{Face: deck.JACK, Suit: deck.CLUB},
			deck.Card{Face: deck.TEN, Suit: deck.CLUB}}}),
		"moving nine onto jack/ten in tableau is ok")
	assert.Error(t, IsMoveLegal(
		tableau,
		deck.Deck{Cards: []deck.Card{deck.Card{Face: deck.JACK, Suit: deck.CLUB}}},
		tableau,
		deck.Deck{Cards: []deck.Card{deck.Card{Face: deck.TEN, Suit: deck.CLUB}}}),
		"moving jack onto ten in tableau is illegal")
}
