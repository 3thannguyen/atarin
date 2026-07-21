package board

import (
	"math/rand"
	"strings"
	"testing"
)

// parseBoard builds a position from an ASCII diagram (row 1 on top).
// '.' empty, 'X' black, 'O' white. sets stones directly then rebuilds
// the chain state, so it does not enforce rules during setup.
func parseBoard(t *testing.T, layout string) *Board {
	t.Helper()
	rows := strings.Split(strings.TrimSpace(layout), "\n")
	size := len(strings.TrimSpace(rows[0]))
	if len(rows) != size {
		t.Fatalf("diagram is %d rows but %d cols wide; want square", len(rows), size)
	}
	b := New(size)
	for r, line := range rows {
		line = strings.TrimSpace(line)
		if len(line) != size {
			t.Fatalf("row %d has width %d, want %d", r+1, len(line), size)
		}
		for c, ch := range line {
			switch ch {
			case 'X':
				b.points[b.Index(r+1, c+1)] = Black
			case 'O':
				b.points[b.Index(r+1, c+1)] = White
			case '.':
				// leave Empty
			default:
				t.Fatalf("unexpected char %q in diagram", ch)
			}
		}
	}
	b.rebuild()
	return b
}

func TestSingleStoneCaptureCenter(t *testing.T) {
	// black surrounds a lone white stone on three sides, then plays the fourth.
	b := parseBoard(t, `
.....
..X..
.XO..
..X..
.....`)
	if !b.Play(b.Index(3, 4), Black) {
		t.Fatal("capturing move should be legal")
	}
	if got := b.colorAt(3, 3); got != Empty {
		t.Fatalf("white stone should be captured, got %v\n%s", got, b)
	}
	if got := b.colorAt(3, 4); got != Black {
		t.Fatalf("black stone should be on the board, got %v\n%s", got, b)
	}
}

func TestCornerCapture(t *testing.T) {
	// corner stone has only two liberties.
	b := parseBoard(t, `
OX...
.....
.....
.....
.....`)
	if !b.Play(b.Index(2, 1), Black) {
		t.Fatal("corner capture should be legal")
	}
	if got := b.colorAt(1, 1); got != Empty {
		t.Fatalf("corner white stone should be captured\n%s", b)
	}
}

func TestMultiStoneCapture(t *testing.T) {
	// a two-stone white chain with a single shared liberty at (3,5).
	b := parseBoard(t, `
.....
..XX.
.XOO.
..XX.
.....`)
	if !b.Play(b.Index(3, 5), Black) {
		t.Fatal("move filling the last liberty should be legal")
	}
	if b.colorAt(3, 3) != Empty || b.colorAt(3, 4) != Empty {
		t.Fatalf("both white stones should be captured\n%s", b)
	}
}

func TestCaptureGivesLibertiesBack(t *testing.T) {
	// After the capture in the multi-stone case, the surrounding black
	// stones must regain liberties where the white stones used to be.
	b := parseBoard(t, `
.....
..XX.
.XOO.
..XX.
.....`)
	b.Play(b.Index(3, 5), Black)
	// The black stone at (2,3) now borders the freshly emptied (3,3).
	if b.libs[b.find(b.Index(2, 3))] == 0 {
		t.Fatalf("surrounding black should have liberties after capture\n%s", b)
	}
}

func TestCapturePriorityOverSuicide(t *testing.T) {
	// white rings the board with its only liberty at the center. black playing
	// the center has zero liberties of its own, but capturing resolves first,
	// so the move is legal and removes all eight white stones.
	b := parseBoard(t, `
OOO
O.O
OOO`)
	if !b.Play(b.Index(2, 2), Black) {
		t.Fatalf("capturing move must beat the suicide rule\n%s", b)
	}
	for r := 1; r <= 3; r++ {
		for c := 1; c <= 3; c++ {
			if r == 2 && c == 2 {
				continue
			}
			if b.colorAt(r, c) != Empty {
				t.Fatalf("all white stones should be captured\n%s", b)
			}
		}
	}
	if b.colorAt(2, 2) != Black {
		t.Fatalf("black stone should remain at center\n%s", b)
	}
}

func TestSuicideIsIllegal(t *testing.T) {
	// white plays into a point fully surrounded by black with no capture.
	b := parseBoard(t, `
.X.
X.X
.X.`)
	if b.Play(b.Index(2, 2), White) {
		t.Fatalf("suicide should be rejected\n%s", b)
	}
	if b.colorAt(2, 2) != Empty {
		t.Fatalf("board must be unchanged after illegal move\n%s", b)
	}
	// The surrounding black stones must be untouched and correct.
	g := b.cloneBoard()
	g.rebuild()
	gInc, lInc := b.groupSummary()
	gReb, lReb := g.groupSummary()
	for i := range gInc {
		if gInc[i] != gReb[i] {
			t.Fatalf("state corrupted by rejected move at %d", i)
		}
	}
	for k, v := range lInc {
		if lReb[k] != v {
			t.Fatalf("liberties corrupted by rejected move")
		}
	}
}

func TestPlayOnOccupiedIsIllegal(t *testing.T) {
	b := parseBoard(t, `
X..
...
...`)
	if b.Play(b.Index(1, 1), White) {
		t.Fatal("playing on an occupied point should fail")
	}
	if b.Play(b.Index(1, 1), Black) {
		t.Fatal("playing on an occupied point should fail even for same color")
	}
}

// TestRandomGameInvariants plays long random games and, after every move,
// checks that the incrementally-maintained chains and pseudo-liberties match
// a from-scratch rebuild.
func TestRandomGameInvariants(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	for game := 0; game < 50; game++ {
		b := New(9)
		var color Color = Black
		passes := 0
		for move := 0; move < 200 && passes < 2; move++ {
			empties := make([]int, 0, 81)
			for r := 1; r <= 9; r++ {
				for c := 1; c <= 9; c++ {
					p := b.Index(r, c)
					if b.points[p] == Empty {
						empties = append(empties, p)
					}
				}
			}
			rng.Shuffle(len(empties), func(i, j int) { empties[i], empties[j] = empties[j], empties[i] })
			played := false
			for _, p := range empties {
				if b.Play(p, color) {
					played = true
					break
				}
			}
			if played {
				passes = 0
			} else {
				passes++
			}

			// compare incremental state against a fresh oracle.
			oracle := b.cloneBoard()
			oracle.rebuild()
			gInc, lInc := b.groupSummary()
			gReb, lReb := oracle.groupSummary()
			for i := range gInc {
				if gInc[i] != gReb[i] {
					t.Fatalf("game %d move %d: chain partition mismatch at point %d\n%s", game, move, i, b)
				}
			}
			if len(lInc) != len(lReb) {
				t.Fatalf("game %d move %d: chain count mismatch (%d vs %d)\n%s", game, move, len(lInc), len(lReb), b)
			}
			for k, v := range lInc {
				if lReb[k] != v {
					t.Fatalf("game %d move %d: pseudo-liberty mismatch for chain %d: incremental=%d oracle=%d\n%s",
						game, move, k, v, lReb[k], b)
				}
			}
			color = color.Opponent()
		}
	}
}
