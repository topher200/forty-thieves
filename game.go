package main

import (
	"errors"
	"fmt"

	"github.com/topher200/baseutil"
	"github.com/topher200/deck"
)

type GameState struct {
	Stock       deck.Deck   `json:"stock"`
	Foundations []deck.Deck `json:"foundations"`
	Tableaus    []deck.Deck `json:"tableaus"`
}

const (
	numTableaus                = 10
	numStartingCardsPerTableau = 4
	numFoundations             = 8
)

// popFromStock returns error if there's no cards in the stock.
func (state *GameState) popFromStock() (deck.Card, error) {
	if len(state.Stock.Cards) <= 0 {
		return deck.Card{}, errors.New("Empty stock")
	}
	card := state.Stock.Cards[0]
	state.Stock.Cards = state.Stock.Cards[1:]
	return card, nil
}

func NewGame() (state GameState) {
	// Combine two decks to make our game deck
	newDeck := deck.NewDeck(false)
	newDeck2 := deck.NewDeck(false)
	newDeck.Cards = append(newDeck.Cards, newDeck2.Cards...)
	newDeck.Shuffle()

	// All cards start in the stock, and our foundations start empty
	state.Stock.Cards = newDeck.Cards
	state.Foundations = make([]deck.Deck, 8)

	// Populate our tableaus with cards off the stock
	state.Tableaus = make([]deck.Deck, 10)
	for i, _ := range state.Tableaus {
		for j := 0; j < numStartingCardsPerTableau; j++ {
			card, err := state.popFromStock()
			baseutil.Check(err)
			state.Tableaus[i].Cards = append(state.Tableaus[i].Cards, card)
		}
	}
	return
}

func (state GameState) String() string {
	str := fmt.Sprintf("stock: %v\n", state.Stock)
	str += "Foundations\n"
	for _, foundation := range state.Foundations {
		str += fmt.Sprintf(" :%v\n", foundation)
	}
	str += "Tableaus\n"
	for _, tableau := range state.Tableaus {
		str += fmt.Sprintf(" :%v\n", tableau)
	}
	return str
}
