package main

import (
	"github.com/topher200/forty-thieves/libgame"
)

// From https://golang.org/pkg/container/heap/ on Sept 3, 2018, modified. We've
// modified the code to adapt to our data structure and to return the lowest
// scoring priority items first.

type PriorityQueue []*libgame.GameState

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Score < pq[j].Score
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*libgame.GameState)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
