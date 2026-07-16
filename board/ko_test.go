package board

import (
	"testing"
)

// we cannot test ko with an ascii diagram since it lacks history -> build with moves
func TestSimpleKoIsIllegal(t *testing.T) {
	b := New(5)
	seq := []struct {
		r, c int
		col  Color
	}{
		{2, 2, Black}, {2, 3, White},
		{3, 1, Black}, {3, 4, White},
		{4, 2, Black}, {4, 3, White},
		{3, 3, Black}, // black pokes into the white shape: the ko stone
	}
	for _, m := range seq {
		if !b.Play(b.Index(m.r, m.c), m.col) {
			t.Fatalf("setup move (%d,%d) should be legal\n%s", m.r, m.c, b)
		}
	}
	if !b.Play(b.Index(3, 2), White) { // the ko capture
		t.Fatalf("white ko capture should be legal\n%s", b)
	}
	if b.Play(b.Index(3, 3), Black) { // immediate recapture: forbidden
		t.Fatalf("immediate ko recapture must be rejected\n%s", b)
	}
	// After a ko-threat exchange elsewhere, the whole-board position
	// differs, so the recapture becomes legal again.
	if !b.Play(b.Index(5, 5), Black) || !b.Play(b.Index(1, 5), White) {
		t.Fatal("ko-threat exchange moves should be legal")
	}
	if !b.Play(b.Index(3, 3), Black) {
		t.Fatalf("ko recapture after an exchange should be legal\n%s", b)
	}
}
