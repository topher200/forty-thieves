package libsolver

import (
	"fmt"

	"github.com/topher200/forty-thieves/libgame"
)

type pile struct {
	location libgame.PileLocation
	index    int
}

func allPiles() []pile {
	piles := make([]pile, 0)
	piles = append(piles, pile{"stock", 0})
	piles = append(piles, pile{"waste", 0})
	for i := 0; i < libgame.NumFoundations; i++ {
		piles = append(piles, pile{"foundation", i})
	}
	for i := 0; i < libgame.NumTableaus; i++ {
		piles = append(piles, pile{"tableau", i})
	}
	return piles
}

func getPossibleMoves(state *libgame.GameState) []libgame.MoveRequest {
	possibleMoves := make([]libgame.MoveRequest, 0)
	piles := allPiles()
	for i, _ := range piles {
		for j, _ := range piles {
			move := libgame.MoveRequest{
				FromPile:  piles[i].location,
				FromIndex: piles[i].index,
				ToPile:    piles[j].location,
				ToIndex:   piles[j].index,
			}
			if state.IsMoveRequestLegal(move) == nil {
				possibleMoves = append(possibleMoves, move)
			}
		}
	}
	return possibleMoves
}

func FoundationAvailableCard(state *libgame.GameState) error {
	for _, move := range getPossibleMoves(state) {
		if move.FromPile != "foundation" && move.ToPile == "foundation" {
			return state.MoveCard(move)
		}
	}
	return fmt.Errorf("No foundationable cards found")
}
