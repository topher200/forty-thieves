package main

import (
	"container/heap"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/topher200/forty-thieves/libgame"
)

// This example creates a PriorityQueue with some items, adds and manipulates an item,
// and then removes the items in priority order.
func TestPriorityQueue(t *testing.T) {
	// Some items and their priorities.
	items := []libgame.GameState{
		libgame.GameState{GameStateID: uuid.NewV4(), Score: 3},
		libgame.GameState{GameStateID: uuid.NewV4(), Score: 2},
		libgame.GameState{GameStateID: uuid.NewV4(), Score: 4},
		libgame.GameState{GameStateID: uuid.NewV4(), Score: 1},
		libgame.GameState{GameStateID: uuid.NewV4(), Score: 3},
		libgame.GameState{GameStateID: uuid.NewV4(), Score: 2},
		libgame.GameState{GameStateID: uuid.NewV4(), Score: 4},
		libgame.GameState{GameStateID: uuid.NewV4(), Score: 1},
		libgame.GameState{GameStateID: uuid.NewV4(), Score: 3},
		libgame.GameState{GameStateID: uuid.NewV4(), Score: 2},
		libgame.GameState{GameStateID: uuid.NewV4(), Score: 4},
		libgame.GameState{GameStateID: uuid.NewV4(), Score: 1},
	}

	// Create a priority queue and put the items in it
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	for _, gs := range items {
		heap.Push(&pq, &gs)
	}

	// Take the items out; they arrive in decreasing priority order.
	previousItem := heap.Pop(&pq).(*libgame.GameState)
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*libgame.GameState)
		assert.Condition(
			t,
			func() bool { return previousItem.Score <= item.Score },
			"%v should score lower than %v", previousItem, item)
		previousItem = item
	}
}
