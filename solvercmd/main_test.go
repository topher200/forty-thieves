package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/topher200/forty-thieves/libgame"
)

// Test that we don't skip good move
func TestShouldSkipMoveDoNotSkipGoodMoves(t *testing.T) {
	assert.False(t, shouldSkipMove(
		libgame.MoveRequest{
			libgame.TABLEAU,
			0,
			libgame.TABLEAU,
			0,
		}))
	assert.False(t, shouldSkipMove(
		libgame.MoveRequest{
			libgame.TABLEAU,
			0,
			libgame.FOUNDATION,
			0,
		}))
	assert.False(t, shouldSkipMove(
		libgame.MoveRequest{
			libgame.STOCK,
			0,
			libgame.TABLEAU,
			0,
		}))
	assert.False(t, shouldSkipMove(
		libgame.MoveRequest{
			libgame.STOCK,
			0,
			libgame.FOUNDATION,
			0,
		}))
	assert.False(t, shouldSkipMove(
		libgame.MoveRequest{
			libgame.STOCK,
			0,
			libgame.WASTE,
			0,
		}))
}

// Test that we skip bad moves
func TestShouldSkipMoveDoSkipBadMoves(t *testing.T) {
	assert.True(t, shouldSkipMove(
		libgame.MoveRequest{
			libgame.FOUNDATION,
			0,
			libgame.FOUNDATION,
			0,
		}))
	assert.True(t, shouldSkipMove(
		libgame.MoveRequest{
			libgame.FOUNDATION,
			0,
			libgame.TABLEAU,
			0,
		}))
}
