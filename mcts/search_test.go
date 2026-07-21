package mcts

import (
	"math/rand/v2"
	"testing"

	"github.com/3thannguyen/atarin/board"
)

func BenchmarkPlayout(bm *testing.B) {
	rng := rand.New(rand.NewPCG(3, 4)) // hardcode permutated congruential generator -> consistent randomness
	root := board.New(9)
	buf := make([]int, 0, 81)
	bm.ResetTimer() // reset time spent on init
	for bm.Loop() { // bm.N runs as many n times as possible (at least 1 full sec) to reduce noise. return ns/operations
		b := root.Clone()
		playout(b, board.Black, 7, rng, buf)
	}
}
