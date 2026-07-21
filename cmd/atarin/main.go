package main

import (
	"flag"
	"os"
	"time"

	"github.com/3thannguyen/atarin/gtp" // for running on Sabaki and play auto matches against other engines
)

func main() {
	budget := flag.Duration("time", 3*time.Second, "thinking time per move")
	workers := flag.Int("workers", 4, "parallel search goroutines")
	flag.Parse()
	gtp.Run(os.Stdin, os.Stdout, gtp.NewEngine(*budget, *workers))
}
