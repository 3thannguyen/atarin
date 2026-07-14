package board

type Color int8 // int8 => 1-byte memory => full board ~100 bytes? => cache-friendly when calling playout

const (
	emptyColor = iota
	Black
	White
	Edge // secret fourth color of border so neighbor loops don't need to check for bounds
)

func (c Color) Opponent() Color {
	switch c {
	case Black:
		return White
	case White:
		return Black
	}
	return c
}

func (c Color) isStone() bool {
	return c == Black || c == White
}

// board is a 1D array with sentinel border of Edges.
// for an NxN board, the size of the array would be (N+2)(N+1)+1; refer to brainstorm doc if forget

type Board struct {
	Size   int
	points []Color // board array

	//using union find to track chains
	parent []int
	stones []int // stone count of chain
	libs   []int

	koPoint int
	hash    uint64
	history map[uint64]bool
}

func New(size int) *Board {
	n := (size+2)*(size+1) + 1
	b := &Board{
		Size:    size,
		points:  make([]Color, n),
		parent:  make([]int, n),
		stones:  make([]int, n),
		libs:    make([]int, n),
		history: make(map[uint64]bool),
		koPoint: -1,
	}

	// filling all the board as an Edge first, then assigning emptyColor spaces after (easier)
	// + setting up union find for chaining
	for i := range b.points {
		b.points[i] = Edge
		b.parent[i] = i
	}
	for r := 1; r < size+1; r++ {
		for c := 1; c < size+1; c++ { // wait... c++ HAHAHAHHAHAHAHAHHAHAH
			b.points[r*(size+1)+c] = emptyColor
		}
	}
	return b
}

// indexing for easier stone placing
func (b *Board) Index(row, col int) int {
	return row*(b.Size+1) + col
}

// using caller-provided buffer to avoid allocating to heap (since playout will have to run this a gazillion times)
func (b *Board) neighbors(p int, buf *[4]int) {
	n := b.Size + 1
	buf[0], buf[1], buf[2], buf[3] = p+1, p-1, p+n, p-n
}

func (b *Board) colorAt(row, col int) Color {
	return b.points[b.Index(row, col)]
}

// Play represents a turn, attempts to place stone at position p
func (b *Board) Play(p int, c Color) bool {
	if !c.isStone() || p < 0 || p > len(b.points) || c != emptyColor {
		return false
	}

	// create a snapshot (copy of board) to replay if illegal move
	// downside is a lot of copies would be created
	snap := b.snapshot()
	b.placeStone(c, p)
	// if liberty of a stone/chain is 0 (an eye = no liberty) then move is illegal
	if b.libs[b.find(p)] == 0 {
		restore(snap)
		return false
	}
	return true
}

func (b *Board) placeStone(p int, c Color) {
	var nb [4]int
	b.neighbors(p, &nb)

	b.parent[p] = p
	b.points[p] = c
	b.stones[p] = 1

	libs := 0
	for _, q := range nb {
		if b.points[q] == emptyColor {
			libs++
		}
	}
	b.libs[p] = libs // setting liberties for stone at p

	for _, q := range nb {
		if b.points[q].isStone() {
			b.libs[b.find(q)]-- // decreasing libs of stones/chains neighbouring p
		}
	}
	
	opp := c.Opponent()
	var enemyRoots []int
	for _, q := range nb{
		if b.points[q] = opp{
			r := b.find(q)
			if b.libs[r] == 0{
				seen := false
				// setting seen so we won't add the same root to enemyRoots
				for _, er := range enemyRoots {
					if er == r{
						seen = true
						break
					}
				}
				if !seen{
					enemyRoots.append(enemyRoots, r)
				}
			}	
		}
	}
	for _, r in := enemyRoots{
		b.removeChain(r, c)	
	}
		
	

}

func (b *Board) snapshot() {

}

func (b *Board) removeChain(root int, friendly Color){
	dead := b.points[root]

}
