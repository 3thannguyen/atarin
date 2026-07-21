package mcts

import (
	"math/rand/v2"

	"github.com/3thannguyen/atarin/board"
)

const maxPlayoutMoves = 300 // safety valve

func playout(b *board.Board, toPlay board.Color, komi float64, rng *rand.Rand, buf []int) board.Color {
	passes := 0
	for move := 0; passes < 2 && move < maxPlayoutMoves; move++ {
		moves := b.CandidateMoves(toPlay, buf)
		played := false

		// we now try candidateMoves in random order until one is legal
		for len(moves) > 0 {
			i := rng.IntN(len(moves))
			m := moves[i]
			// swap-remove instead of only removing O(1) prevents shifting everything in the slice
			moves[i] = moves[len(moves)-1]
			moves = moves[:len(moves)-1]
			if b.Play(m, toPlay) { // legal check
				played = true
				break
			}
		}
		if played {
			passes = 0
		} else {
			passes++
		}
		toPlay = toPlay.Opponent()
	}
	// determining who wins (since Score() returns difference in scores between Black and White), black must win by more than komi
	if float64(b.Score()) > komi {
		return board.Black
	}
	return board.White
}
