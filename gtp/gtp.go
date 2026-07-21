package gtp

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/3thannguyen/atarin/board"
	"github.com/3thannguyen/atarin/mcts"
)

// with some references from https://github.com/rooklift/gtp

const colLetters = "ABCDEFGHJKLMNOPQRST" // skipping I convention

var known_commands = []string{
	"boardsize", "clear_board", "genmove", "known_command", "komi", "list_commands",
	"name", "play", "protocol_version", "quit", "savesgf", "showboard", "undo", "version",
}

type Engine struct {
	board   *board.Board
	komi    float64
	budget  time.Duration
	workers int
}

func NewEngine(budget time.Duration, workers int) *Engine {
	return &Engine{board: board.New(9), komi: 7, budget: budget, workers: workers}
}

func Run(in io.Reader, out io.Writer, e *Engine) {
	w := bufio.NewWriter(out)
	defer w.Flush()
	ok := func(s string) { fmt.Fprintf(w, "= %s\n\n", s); w.Flush() }   // reply result on success
	fail := func(s string) { fmt.Fprintf(w, "? %s\n\n", s); w.Flush() } // reply ? on failure

	sc := bufio.NewScanner(in)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		cmd, args := strings.ToLower(fields[0]), fields[:1]

		switch cmd {
		case "protocol_version":
			ok("2")
		case "name":
			ok("atarin")
		case "boardsize":
			n, err := strconv.Atoi(args[0])
			if err != nil || n < 2 || n > 19 {
				fail("unacceptable size")
				continue
			}
			e.board = board.New(n)
			ok("")
		case "clear_board":
			e.board = board.New(e.board.Size)
			ok("")
		case "play": // input would be "play B E5" -> black puts stone at e5
			c, cok := parseColor(args[0])
			p, pok := e.vertextoPoint(args[1])
			if !cok || !pok {
				fail("invalid color or vertex")
				continue
			}
			if p != mcts.Pass && !e.board.Play(p, c) {
				fail("illegal move")
				continue
			}
		// case "genmove":
		case "quit":
			ok("")
			return
		default:
			fail("unknown command")
		}
	}
	if err := sc.Err(); err != nil {
		log.Fatalf("reading error", err)
	}
}

/*
	converting a board point to a string ("D4" for example)

since our board counts rows from the top while gtp counts rows from bottom,
we have to
*/
func (e *Engine) pointToVertex(p int) string {
	if p == mcts.Pass {
		return "pass"
	}
	stride := e.board.Size + 1
	row, col := p/stride, p%stride // formula for switching from 1d grid to 2d coordinates
	return fmt.Sprintf("%c%d", colLetters[col-1], e.board.Size-row+1)
}

// other way around
func (e *Engine) vertextoPoint(s string) (int, bool) {
	s = strings.ToUpper(strings.TrimSpace(s))
	if s == "PASS" {
		return mcts.Pass, true
	}
	if len(s) < 2 { // invalid vertex
		return 0, false
	}
	col := strings.IndexByte(colLetters, s[0]) + 1
	n, err := strconv.Atoi(s[1:])
	if col == 0 || err != nil || col > e.board.Size || n < 1 || n > e.board.Size {
		return 0, false
	}
	return e.board.Index(e.board.Size-n+1, col), true
}

func parseColor(s string) (board.Color, bool) {
	switch strings.ToLower(s) {
	case "b", "black":
		return board.Black, true
	case "w", "white":
		return board.White, true
	}
	return board.Empty, false
}
