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
