package deck

import "fmt"

// Deck is a deck of cards. An array of type Card
type Deck struct {
	Cards []Card
}

func (d Deck) String() string {
	str := ""
	for _, card := range d.Cards {
		str += fmt.Sprint(card) + " "
	}
	return str
}

// NewDeck creates and returns a new deck with the bool parameter to either shuffle (true) or non shuffle (false)
func NewDeck(shuffled bool) Deck {
	deck := NewSpecificDeck(shuffled, FACES, SUITS)
	return deck
}

// NewSpecificDeck creates and returns a deck that is created
// with all the premutations of an array of Faces and an array of Suits.
// The same bool parameter is expected to shuffle the deck
func NewSpecificDeck(shuffled bool, faces []Face, suits []Suit) Deck {
	cards := make([]Card, len(suits)*len(faces))
	for sindex, suit := range suits {
		for findex, face := range faces {
			index := (sindex * len(faces)) + findex
			cards[index] = Card{face, suit}
		}
	}
	deck := Deck{cards}
	if shuffled {
		deck.Shuffle()
	}
	return deck
}

// NewEmptyDeck creates an empty deck with an empty array of Cards
func NewEmptyDeck() Deck {
	deck := Deck{[]Card{}}
	return deck
}

// NumberOfCards is a utility function that tells you how many cards are left in the deck
func (d *Deck) NumberOfCards() int {
	return len(d.Cards)
}
