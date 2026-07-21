package board

// root finding logic (union find) in order to simplify keeping metadata of each stone
// finding the root of chain p (using path halving), union find method
func (b *Board) find(p int) int {
	for b.parent[p] != p {
		b.parent[p] = b.parent[b.parent[p]]
		p = b.parent[p]
	}
	return p
}

func (b *Board) union(p, q int) {
	rp, rq := b.find(p), b.find(q)
	if rp == rq {
		return
	}
	if b.stones[rp] < b.stones[rq] {
		rp, rq = rq, rp
	}
	// updating parents + chain 'length' and liberties
	b.parent[rq] = rp
	b.stones[rp] += b.stones[rq]
	b.libs[rp] += b.libs[rq]
}

// rebuild recomputes entire union-find and liberty state from raw stones
// this will be used for a slow, obviously correct testing against fast one
func (b *Board) rebuild() {
	for i := range b.parent {
		b.parent[i] = i
		b.stones[i] = 0
		b.libs[i] = 0
	}
	visited := make([]bool, len(b.points))
	for p := range b.points {
		col := b.points[p]
		if !col.isStone() || visited[p] {
			continue
		}
		visited[p] = true
		stack := []int{p}
		var members []int
		for len(stack) > 0 {
			s := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			members = append(members, s)
			var nb [4]int
			b.neighbors(s, &nb)
			for _, q := range nb {
				if b.points[q] == col && !visited[q] {
					visited[q] = true
					stack = append(stack, q)
				}
			}
		}
		root := members[0]
		libs := 0
		for _, s := range members {
			b.parent[s] = root
			var nb [4]int
			b.neighbors(s, &nb)
			for _, q := range nb {
				if b.points[q] == Empty {
					libs++
				}
			}
		}
		b.stones[root] = len(members)
		b.libs[root] = libs
	}
	b.hash = b.recomputeHash()
	if b.history == nil {
		b.history = make(map[uint64]bool)
	}
	b.history[b.hash] = true
}

// returns deep-copy so we can build a from-scratch oracle test without disturbing original board
func (b *Board) cloneBoard() *Board {
	hist := make(map[uint64]bool, len(b.history))
	for k := range b.history {
		hist[k] = true
	}
	return &Board{
		Size:    b.Size,
		points:  append([]Color(nil), b.points...),
		parent:  append([]int(nil), b.parent...),
		stones:  append([]int(nil), b.stones...),
		libs:    append([]int(nil), b.libs...),
		hash:    b.hash,
		history: hist,
	}
}

/* groupSummary() return the smallest id of its union-find chain as well as
the map from that chain id to its pseudo-liberty count. basically to canonicalize
the group arrangements (regardless of union order). check if chains are actually tracked
correctly
*/

func (b *Board) groupSummary() (groupOf []int, libsByGroup map[int]int) {
	groupOf = make([]int, len(b.points))
	for i := range groupOf {
		groupOf[i] = -1
	}
	rootMin := map[int]int{}
	for p := range b.points {
		if b.points[p].isStone() {
			r := b.find(p)
			if m, ok := rootMin[r]; !ok || p < m {
				rootMin[r] = p
			}
		}
	}
	for p := range b.points {
		if b.points[p].isStone() {
			groupOf[p] = rootMin[b.find(p)]
		}
	}
	libsByGroup = map[int]int{}
	for r, m := range rootMin {
		libsByGroup[m] = b.libs[r]
	}
	return groupOf, libsByGroup
}
