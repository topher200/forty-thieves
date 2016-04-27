package deck

import "fmt"

// Suit represents the suit of the card (spade, heart, diamon, club)
type Suit string

// Face represents the face of the card (ace, two...queen, king)
type Face string

// Contants for Suit ♠♥♦♣
const (
	CLUB    Suit = "clubs"
	DIAMOND      = "diamonds"
	HEART        = "hearts"
	SPADE        = "spades"
)

// Contants for Face
const (
	ACE   Face = "A"
	TWO        = "2"
	THREE      = "3"
	FOUR       = "4"
	FIVE       = "5"
	SIX        = "6"
	SEVEN      = "7"
	EIGHT      = "8"
	NINE       = "9"
	TEN        = "T"
	JACK       = "J"
	QUEEN      = "Q"
	KING       = "K"
)

// findIndex returns the location in the FACES slice of face
func (face Face) findIndex() (int, error) {
	for k, v := range FACES {
		if face == v {
			return k, nil
		}
	}
	return 0, fmt.Errorf("face not found in faces list: '%s'", face)
}

// Decrement subtracts one value from a Face.
//
// Ace is considered low - decrementing from it returns error
func Decrement(face Face) (Face, error) {
	index, err := face.findIndex()
	if err != nil {
		return face, fmt.Errorf("Face '%v' not found", face)
	}
	if index <= 0 {
		return face, fmt.Errorf("Can't decrement lowest Face '%s'", face)
	}
	return FACES[index-1], nil
}

// Increment adds one value to a Face.
//
// King is considered high - incrementing from it returns error
func Increment(face Face) (Face, error) {
	index, err := face.findIndex()
	if err != nil {
		return face, fmt.Errorf("Face '%v' not found", face)
	}
	if index >= len(FACES)-1 {
		return face, fmt.Errorf("Can't increment highest Face '%s'", face)
	}
	return FACES[index+1], nil
}

// Global Variables representing the default suits and faces in a deck of cards
var (
	SUITS = []Suit{CLUB, DIAMOND, HEART, SPADE}
	FACES = []Face{ACE, TWO, THREE, FOUR, FIVE, SIX, SEVEN, EIGHT, NINE, TEN, JACK, QUEEN, KING}
)
