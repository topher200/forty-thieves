package libgame

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/topher200/baseutil"
	"github.com/topher200/deck"
)

type Game struct {
	ID int64
}

// TODO(topher): this probably shouldn't have an ID field
type GameState struct {
	GameID            int64
	GameStateID       uuid.UUID
	PreviousGameState uuid.NullUUID
	MoveNum           int64
	Stock             deck.Deck
	Foundations       []deck.Deck
	Tableaus          []deck.Deck
	Waste             deck.Deck
	Score             int // Must be updated after any modifications to the Decks above
}

const (
	NumFoundations             = 8
	NumTableaus                = 10
	numStartingCardsPerTableau = 4
)

// popFromStock returns error if there's no cards in the stock.
//
// Doesn't call 'updateScore' because it's a private function. Caller is
// expected to do that for us.
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

// IsMoveRequestLegal requests legality of a MoveRequest for a given GameState.
//
// Returns an error (with explanation) if move shouldn't be done
func (state *GameState) IsMoveRequestLegal(move MoveRequest) error {
	fromDeck, toDeck, err := state.parseDecksFromMoveRequest(move)
	if err != nil {
		return err
	}

	return isMoveLegal(move.FromPile, fromDeck, move.ToPile, toDeck)
}

// isMoveLegal checks the cards and decks involved for legality.
func isMoveLegal(
	fromPile PileLocation, fromDeck *deck.Deck,
	toPile PileLocation, toDeck *deck.Deck) error {

	// Is the origin location illegal?
	if fromPile == stock {
		return fmt.Errorf("Illegal move - origin '%s' illegal", fromPile)
	}

	// Is the destination location illegal?
	if toPile == stock || toPile == waste {
		return fmt.Errorf("Illegal move - destination '%s' illegal", toPile)
	}

	// Is there a card to move?
	if len(fromDeck.Cards) <= 0 {
		return fmt.Errorf("Illegal move - 'from' pile '%s' empty", fromPile)
	}
	cardBeingMoved := fromDeck.Cards[len(fromDeck.Cards)-1]

	// Are our destination empty?
	if len(toDeck.Cards) <= 0 {
		// Empty foundations can only take aces
		if toPile == foundation && cardBeingMoved.Face != deck.ACE {
			return fmt.Errorf(
				"Illegal move - moving to empty foundation requires ACE, not '%s'",
				cardBeingMoved)
		}
		// Empty tableaus are always OK moves
	} else {
		destinationCard := toDeck.Cards[len(toDeck.Cards)-1]
		if cardBeingMoved.Suit != destinationCard.Suit {
			return fmt.Errorf("Illegal move - suits much match (%s on %s)",
				cardBeingMoved, destinationCard)
		}
		switch toPile {
		case tableau:
			decrementedDestination, err := deck.Decrement(destinationCard.Face)
			if err != nil {
				return err
			}
			if decrementedDestination != cardBeingMoved.Face {
				return fmt.Errorf("Illegal move - tableau cards must decrease (%s on %s)",
					cardBeingMoved, destinationCard)
			}
		case foundation:
			incrementedDestination, err := deck.Increment(destinationCard.Face)
			if err != nil {
				return err
			}
			if incrementedDestination != cardBeingMoved.Face {
				return fmt.Errorf("Illegal move - foundation cards must increase (%s on %s)",
					cardBeingMoved, destinationCard)
			}
		}
	}
	return nil
}

type MoveRequest struct {
	FromPile  PileLocation
	FromIndex int
	ToPile    PileLocation
	ToIndex   int
}

func (state *GameState) parseDecksFromMoveRequest(
	move MoveRequest) (*deck.Deck, *deck.Deck, error) {
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

func (state *GameState) MoveCard(move MoveRequest) error {
	err := state.IsMoveRequestLegal(move)
	if err != nil {
		return fmt.Errorf("Can't complete move: %v", err)
	}
	fromDeck, toDeck, err := state.parseDecksFromMoveRequest(move)
	if err != nil {
		return err
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

// DealNewGame takes a game and randomly deals a starting gamestate for that game
func DealNewGame(game Game) (state GameState) {
	state.GameStateID = uuid.NewV4()
	state.GameID = game.ID
	state.MoveNum = 0

	// Combine two decks to make our game deck
	newDeck := deck.NewDeck(false)
	newDeck2 := deck.NewDeck(false)
	newDeck.Cards = append(newDeck.Cards, newDeck2.Cards...)
	newDeck.Shuffle()

	// All cards start in the stock, and our foundations start empty
	state.Stock.Cards = newDeck.Cards
	state.Foundations = make([]deck.Deck, NumFoundations)

	// Populate our tableaus with cards off the stock
	state.Tableaus = make([]deck.Deck, NumTableaus)
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
