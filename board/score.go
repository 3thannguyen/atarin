package board

/* counting region right now goes with tromp-taylor scoring:
score = your stones + your area
if the empty space only touch one color then it is that color's region
if touches both then belongs to none. this can be done by two boolean expressions.
*/

func (b *Board) Score() int {
	black, white := 0, 0
	seen := make([]bool, len(b.points))
	for p, c := range b.points {
		switch c {
		case Black:
			black++
		case White:
			white++
		case Empty:
			// checking all empty stones of 1 region and their neighbors
			if seen[p] {
				continue
			}
			seen[p] = true
			region := 0
			touchB, touchW := false, false
			stack := []int{p}
			for len(stack) > 0 {
				region++
				s := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				var nb [4]int
				b.neighbors(s, &nb)
				for _, q := range nb {
					switch b.points[q] {
					case Empty:
						if !seen[q] {
							seen[q] = true
							stack = append(stack, q)
						}
					case Black:
						touchB = true
					case White:
						touchW = true
					}
				}
			}
			if touchB && !touchW {
				black += region
			} else if touchW && !touchB {
				white += region
			}
		}
	}
	// returns score in this format before komi. komi can be configured later
	return black - white
}

/* an important aspect of go is checking if an eye is real or false;
i still have yet to think of a way to properly implement it, but apparently a practical rule
is that an eye is real if the opponent has <= 1 stone in the diags if the eye is in center,
none if the eye is on the edge. having a check on whether we have a simply eye would be helpful
for the agent's playmaking as well (so that it doesn't play in its own eye)
*/

func (b *Board) isSimplyEye(p int, c Color) bool {
	var nb [4]int
	b.neighbors(p, &nb)
	// checking if all neighbors are our color/edge
	for _, q := range nb {
		if b.points[q] != c && b.points[q] != Edge {
			return false
		}
	}
	n := b.Size + 1
	diags := [4]int{p - n - 1, p - n + 1, p + n - 1, p + n + 1}
	enemy, edge := 0, 0
	opp := c.Opponent()
	for _, d := range diags {
		switch b.points[d] {
		case opp:
			enemy++
		case Edge:
			edge++
		}
	}
	if edge > 0 {
		return enemy == 0
	}
	return enemy <= 1
}

// giving a list of moves for mcts to play (cheap and optimistic), integrates isSimpleEye()
func (b *Board) CandidateMoves(c Color, buf []int) []int {
	buf = buf[:0] // emptying the buffer array but keeping capacity
	for r := 1; r <= b.Size; r++ {
		base := r + (b.Size + 1)
		for col := 1; col <= b.Size; col++ {
			p := base + col
			if b.points[p] == Empty && !b.isSimplyEye(p, c) {
				buf = append(buf, p)
			}
		}
	}
	return buf
}

func (b *Board) Clone() *Board { return b.cloneBoard() }
