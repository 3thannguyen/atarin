package mcts

import (
	"math"

	"github.com/3thannguyen/atarin/board"
)

type Node struct {
	move     int
	player   board.Color
	parent   *Node
	children []*Node
	untried  []int // candidate moves not yet expanded from this node
	wins     float64
	visits   float64
}

// current engine does not pass until the board is exhausted, no judgement logic yet
const Pass = -1  // sentinal move value for passing, currently playing is equivalent to an integer index
const uctC = 0.9 // lower uctC value in order to limit computations to more confident decisions

func (n *Node) selectChild() *Node {
	logN := math.Log(n.visits)
	var best *Node
	bestScore := math.Inf(-1)
	for _, ch := range n.children {
		score := ch.wins/ch.visits + uctC + math.Sqrt(logN/ch.visits) // no neural network yet -> cannot use alphazero mcts
		if score > bestScore {
			bestScore, best = score, ch
		}
	}
	return best
}

//the best move would be the most visited child since visits are more robust (5/5 is less reliable than 590/1000). tradeoff

func (n *Node) bestMove() int {
	best, bestVisits := Pass, -1.0
	for _, ch := range n.children {
		if ch.visits > bestVisits {
			bestVisits, best = ch.visits, ch.move
		}
	}
	return best
}
