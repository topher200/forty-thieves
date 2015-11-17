package main

import (
	"errors"
	"fmt"

	"github.com/topher200/baseutil"
	"github.com/topher200/deck"
)

type pile struct {
	cards []deck.Card
}

type gameState struct {
	stock       pile
	foundations []pile
	tableaus    []pile
}

const (
	numTableaus                = 10
	numStartingCardsPerTableau = 4
	numFoundations             = 8
)

// popFromStock returns error if there's no cards in the stock.
func (state gameState) popFromStock() (deck.Card, error) {
	if len(state.stock.cards) <= 0 {
		return deck.Card{}, errors.New("Empty stock")
	}
	card := state.stock.cards[0]
	state.stock.cards = state.stock.cards[1:]
	return card, nil
}

func NewGame() (state gameState) {
	newDeck := deck.NewDeck(false)
	newDeck2 := deck.NewDeck(false)
	newDeck.Cards = append(newDeck.Cards, newDeck2.Cards...)
	fmt.Println(newDeck)
	state.stock.cards = newDeck.Cards
	state.foundations = make([]pile, 8)

	state.tableaus = make([]pile, 10)
	for _, tableau := range state.tableaus {
		for i := 0; i < numStartingCardsPerTableau; i++ {
			card, err := state.popFromStock()
			baseutil.Check(err)
			tableau.cards = append(tableau.cards, card)
		}
	}
	return
}
