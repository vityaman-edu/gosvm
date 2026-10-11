// GoSVM is a symbolic virtual machine command-line application.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/vityaman-edu/gosvm/internal/gosmv"
	"github.com/vityaman-edu/gosvm/internal/gossa"
	"github.com/vityaman-edu/gosvm/internal/srcloc"
)

func main() {
	err := run(os.Args[1:], os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gosvm:", err)
		os.Exit(1)
	}
}

func run(arguments []string, stdin io.Reader, stdout io.Writer) error {
	options, err := parseArgs(arguments)
	if errors.Is(err, flag.ErrHelp) {
		return writeUsage(stdout)
	}

	if err != nil {
		return err
	}

	source, err := readSource(options.Program, stdin)
	if err != nil {
		return err
	}

	analysis, err := gosmv.Analyze(source, options.Function)
	if err != nil {
		return fmt.Errorf("analyze program: %w", err)
	}

	if !options.PrintSSA {
		return nil
	}

	err = gossa.PrintSSA(stdout, analysis.Package)
	if err != nil {
		return fmt.Errorf("print SSA: %w", err)
	}

	return nil
}

func writeUsage(stdout io.Writer) error {
	_, err := io.WriteString(stdout, usage)
	if err != nil {
		return fmt.Errorf("write usage: %w", err)
	}

	return nil
}

func readSource(path string, stdin io.Reader) (srcloc.Source, error) {
	var text []byte

	var err error

	if path == stdinPath {
		text, err = io.ReadAll(stdin)
	} else {
		text, err = os.ReadFile(filepath.Clean(path))
	}

	if err != nil {
		return srcloc.Source{}, fmt.Errorf("read program %q: %w", path, err)
	}

	return srcloc.Source{Path: path, Text: string(text)}, nil
}
