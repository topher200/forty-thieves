package deck

import "math/rand"

// MultiShuffle calls Shuffle multipule times
func (deck *Deck) MultiShuffle(iterations int) {
	for i := 0; i < iterations; i++ {
		deck.Shuffle()
	}
}

// Shuffle uses Knuth shuffle algo to randomize the deck in O(n) time
// sourced from https://gist.github.com/quux00/8258425
func (deck *Deck) Shuffle() {
	N := len(deck.Cards)
	for i := 0; i < N; i++ {
		r := i + rand.Intn(N-i)
		deck.Cards[r], deck.Cards[i] = deck.Cards[i], deck.Cards[r]
	}
}
