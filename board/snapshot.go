package board

// keeping a boardState so that we can copy/undo
type boardState struct {
	points []Color
	parent []int
	stones []int
	libs   []int
	hash   uint64
}

// taking a snapshot of current board by adding onto a boardState struct
func (b *Board) snapshot() boardState {
	return boardState{
		points: append([]Color(nil), b.points...),
		parent: append([]int(nil), b.parent...),
		stones: append([]int(nil), b.stones...),
		libs:   append([]int(nil), b.libs...),
		hash:   b.hash,
	}
}

// restoring the current board to before the illegal play was made
func (b *Board) restore(s boardState) {
	copy(b.points, s.points)
	copy(b.parent, s.parent)
	copy(b.stones, s.stones)
	copy(b.libs, s.libs)
	b.hash = s.hash
}
