package main

import (
	"flag"
	"io"
	"os"

	"github.com/vityaman-edu/gosvm/internal/gosvm"
)

func ParseText() string {
	text, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	return string(text)
}

func main() {
	function := flag.String("function", "", "function name to analyze")
	printSSA := flag.Bool("print-ssa", false, "print SSA to stderr")
	printSMT := flag.Bool("print-smt", false, "print SMT to stderr")

	flag.Parse()

	gosvm.Go(gosvm.Input{
		Text:     ParseText(),
		Function: *function,
		PrintSSA: *printSSA,
		PrintSMT: *printSMT,
	})
}
