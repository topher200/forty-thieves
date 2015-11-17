package main

import (
	"errors"
	"fmt"

	"github.com/topher200/baseutil"
	"github.com/topher200/deck"
)

type gameState struct {
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
func (state gameState) popFromStock() (deck.Card, error) {
	if len(state.stock.Cards) <= 0 {
		return deck.Card{}, errors.New("Empty stock")
	}
	card := state.stock.Cards[0]
	state.stock.Cards = state.stock.Cards[1:]
	return card, nil
}

func NewGame() (state gameState) {
	newDeck := deck.NewDeck(false)
	newDeck2 := deck.NewDeck(false)
	newDeck.Cards = append(newDeck.Cards, newDeck2.Cards...)
	fmt.Println(newDeck)
	state.stock.Cards = newDeck.Cards
	state.foundations = make([]deck.Deck, 8)

	state.tableaus = make([]deck.Deck, 10)
	for _, tableau := range state.tableaus {
		for i := 0; i < numStartingCardsPerTableau; i++ {
			card, err := state.popFromStock()
			baseutil.Check(err)
			tableau.Cards = append(tableau.Cards, card)
		}
	}
	return
}
