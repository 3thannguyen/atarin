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
