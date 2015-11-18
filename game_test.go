package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/topher200/deck"
)

func TestPopFromStock(t *testing.T) {
	state := NewGame()
	card, err := state.popFromStock()
	assert.Nil(t, err)
	assert.NotEqual(t, card, deck.Card{})
}
