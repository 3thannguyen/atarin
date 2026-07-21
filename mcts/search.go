package mcts

import (
	"math/rand/v2"
	"time"

	"github.com/3thannguyen/atarin/board"
)

/*
four steps that mcts will follow:
1. selection: walk from root down the path of most promising leaves chosen. uses uct to aid.
2. expansion: adds one new child to the tree
3. simulation: run a super quick playout until the end to see who wins
4. back-propagation: travel back the path, updating each nodes: visits += 1 and wins += 1 for the side that won

what we will do on the board is first cloning the root for the mcts, uct walk down the known tree,
replaying each edge's move on clone so board can track the tree path (replaying instead of keeping a board
copy in every node to save space -> position of board is derivable from each step), we then pop a random untried move
and expand one child -> run quick simulation, then back-propagation

make sure that selection's plays are legal (do this by making sure b.Play is called when expanding)
*/

func search(root *board.Board, toPlay board.Color, komi float64, deadline time.Time, rng *rand.Rand) *Node {
	rootNode := &Node{
		move:    Pass, // assuming the opponent passed -> current board position + our turn
		player:  toPlay.Opponent(),
		untried: root.CandidateMoves(toPlay, nil),
	}
	buf := make([]int, 0, root.Size*root.Size) // playout scratch, reused

	for time.Now().Before(deadline) {
		b := root.Clone() // some function clone
		node := rootNode
		color := toPlay

		// selection
		for len(node.untried) == 0 && len(node.children) > 0 {
			node = node.selectChild()
			b.Play(node.move, node.player)
			color = node.player.Opponent()
		}

		// expansion
		for len(node.untried) > 0 {
			i := rng.IntN(len(node.untried))
			m := node.untried[i]
			node.untried[i] = node.untried[len(node.untried)-1]
			node.untried = node.untried[:len(node.untried)-1]
			if !b.Play(m, color) {
				continue
			}
			child := &Node{
				move:    m,
				player:  color,
				parent:  node,
				untried: b.CandidateMoves(color.Opponent(), nil),
			}
			node.children = append(node.children, child)
			node = child
			color = color.Opponent()
			break
		}

		// simulation
		winner := playout(b, color, komi, rng, buf)

		// backpropagation
		for n := node; n != nil; n = n.parent {
			n.visits++
			if winner == n.player {
				n.wins++
			}
		}
	}
	return rootNode
}
