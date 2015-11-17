package main

import (
	"errors"
	"fmt"

	"github.com/topher200/baseutil"
	"github.com/topher200/deck"
)

type GameState struct {
	stock       deck.Deck
	foundations []deck.Deck
	tableaus    []deck.Deck
}

const (
	numTableaus                = 10
	numStartingCardsPerTableau = 4
	numFoundations             = 8
)

// popFromStock returns error if there's no cards in the stock.
func (state *GameState) popFromStock() (deck.Card, error) {
	if len(state.stock.Cards) <= 0 {
		return deck.Card{}, errors.New("Empty stock")
	}
	card := state.stock.Cards[0]
	state.stock.Cards = state.stock.Cards[1:]
	return card, nil
}

func NewGame() (state GameState) {
	// Combine two decks to make our game deck
	newDeck := deck.NewDeck(false)
	newDeck2 := deck.NewDeck(false)
	newDeck.Cards = append(newDeck.Cards, newDeck2.Cards...)
	newDeck.Shuffle()

	// All cards start in the stock, and our foundations start empty
	state.stock.Cards = newDeck.Cards
	state.foundations = make([]deck.Deck, 8)

	// Populate our tableaus with cards off the stock
	state.tableaus = make([]deck.Deck, 10)
	for i, _ := range state.tableaus {
		for j := 0; j < numStartingCardsPerTableau; j++ {
			card, err := state.popFromStock()
			baseutil.Check(err)
			state.tableaus[i].Cards = append(state.tableaus[i].Cards, card)
		}
	}
	return
}

func (state GameState) String() string {
	str := fmt.Sprintf("stock: %v\n", state.stock)
	str += "Foundations\n"
	for _, foundation := range state.foundations {
		str += fmt.Sprintf(" :%v\n", foundation)
	}
	str += "Tableaus\n"
	for _, tableau := range state.tableaus {
		str += fmt.Sprintf(" :%v\n", tableau)
	}
	return str
}
