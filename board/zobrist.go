package board

import (
	"math/rand"
)

/*
zobrist hashing relies entirely on bitwise xor -> adding and removing a stone
uses the same xor. also extremely efficient at looking at previous board positions
by assigning each point a random 64-bit key, then xor of keys of all the stones
on the board would be the board's position hash

that said, there may be collisions (aka positions that map to the same hash), but
since there are so many different positions and so many keys, this is functionally negligible

the tracking of zobrist hashing will allow us to enforce superko
*/

// setting the 1d board layout
const maxPoints = (19+2)*(19+1) + 1

var zobristKeys [maxPoints][2]uint64 // zobristKeys[p][0] is black at p, 1 is white

func init() {
	rng := rand.New(rand.NewSource(0x60BA9))
	// setting a unique value for each color at each point
	for i := range maxPoints {
		zobristKeys[i][0] = rng.Uint64()
		zobristKeys[i][1] = rng.Uint64()
	}
}

func colorSlot(c Color) int {
	if c == Black {
		return 0
	} else {
		return 1
	}
}

func (b *Board) recomputeHash() uint64 {
	var h uint64
	for p, c := range b.points {
		if c.isStone() {
			h ^= zobristKeys[p][colorSlot(c)]
		}
	}
	return h
}
