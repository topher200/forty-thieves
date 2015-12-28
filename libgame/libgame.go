package libgame

import (
	"errors"
	"fmt"

	"github.com/topher200/baseutil"
	"github.com/topher200/deck"
)

type GameState struct {
	Stock       deck.Deck
	Foundations []deck.Deck
	Tableaus    []deck.Deck
	Waste       deck.Deck
	// Must be updated after any modifications to the Decks above
	Score int
}

const (
	numTableaus                = 10
	numStartingCardsPerTableau = 4
	numFoundations             = 8
)

// popFromStock returns error if there's no cards in the stock.
//
// Doesn't call 'updateScore' because it's a private function.
func (state *GameState) popFromStock() (deck.Card, error) {
	if len(state.Stock.Cards) <= 0 {
		return deck.Card{}, errors.New("Empty stock")
	}
	card := state.Stock.Cards[0]
	state.Stock.Cards = state.Stock.Cards[1:]
	return card, nil
}

// Pile locations
type PileLocation string

const (
	stock      PileLocation = "stock"
	foundation              = "foundation"
	tableau                 = "tableau"
	waste                   = "waste"
)

type Move struct {
	FromPile  PileLocation
	FromIndex int
	ToPile    PileLocation
	ToIndex   int
}

func (state *GameState) MoveIsLegal(move Move) bool {
	// todo
	return false
}

func (state *GameState) parseMove(move Move) (*deck.Deck, *deck.Deck, error) {
	parseFunc := func(pileLocation PileLocation, index int) (*deck.Deck, error) {
		var d *deck.Deck
		switch pileLocation {
		case "tableau":
			d = &state.Tableaus[index]
		case "foundation":
			d = &state.Foundations[index]
		case "stock":
			d = &state.Stock
		case "waste":
			d = &state.Waste
		default:
			return nil, fmt.Errorf("unknown pile name '%s'", pileLocation)
		}
		return d, nil
	}
	from, err := parseFunc(move.FromPile, move.FromIndex)
	if err != nil {
		return nil, nil, err
	}
	to, err := parseFunc(move.ToPile, move.ToIndex)
	if err != nil {
		return nil, nil, err
	}
	return from, to, nil
}

func (state *GameState) MoveCard(move Move) error {
	fromDeck, toDeck, err := state.parseMove(move)
	if err != nil {
		return fmt.Errorf("Can't parse Move: %v", err)
	}

	if len(fromDeck.Cards) <= 0 {
		return errors.New("From deck empty - can't complete move")
	}
	toDeck.Cards = append(toDeck.Cards, fromDeck.Cards[len(fromDeck.Cards)-1])
	fromDeck.Cards = fromDeck.Cards[:len(fromDeck.Cards)-1]

	state.updateScore()
	return nil
}

func (state *GameState) FlipStock() error {
	card, err := state.popFromStock()
	if err != nil {
		return errors.New("Can't flip empty stock")
	}

	state.Waste.Cards = append(state.Waste.Cards, card)

	state.updateScore()
	return nil
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

	state.updateScore()
	return
}

func (state GameState) String() string {
	str := fmt.Sprintf("Stock: %v\n", state.Stock)
	str += "Foundations\n"
	for _, foundation := range state.Foundations {
		str += fmt.Sprintf(" :%v\n", foundation)
	}
	str += "Tableaus\n"
	for _, tableau := range state.Tableaus {
		str += fmt.Sprintf(" :%v\n", tableau)
	}
	str += fmt.Sprintf("Waste: %v\n", state.Waste)
	return str
}

// a GameState's Score is the number of cards not in foundations.
//
// The game is won when score is 0.
//
// This function must be called after any function that manipulates the Decks.
func (state *GameState) updateScore() {
	score := 0
	score += len(state.Stock.Cards)
	for i := range state.Tableaus {
		score += len(state.Tableaus[i].Cards)
	}
	score += len(state.Waste.Cards)
	state.Score = score
}
