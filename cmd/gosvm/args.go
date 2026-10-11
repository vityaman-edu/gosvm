package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/vityaman-edu/gosvm/internal/srcloc"
)

const stdinPath = "/dev/stdin"

const usage = "Usage: gosvm [--program=<path>] [--function=<name>] " +
	"[--print-ssa] [program [function]]\n"

var (
	errUnexpectedPositionals = errors.New("unexpected positional arguments")
	errEmptyProgram          = errors.New("program path must not be empty")
	errEmptyFunction         = errors.New("function name must not be empty")
)

type args struct {
	Program  string
	Function srcloc.FunctionName
	PrintSSA bool
}

type argumentGroups struct {
	flags      []string
	positional []string
}

func parseArgs(arguments []string) (args, error) {
	options := args{Program: stdinPath, Function: "main", PrintSSA: false}
	function := string(options.Function)
	flags := flag.NewFlagSet("gosvm", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&options.Program, "program", options.Program,
		"Go source path (defaults to stdin)")
	flags.StringVar(&function, "function", function, "function to analyze")
	flags.BoolVar(&options.PrintSSA, "print-ssa", false,
		"print package SSA to stdout")

	groups := splitArguments(arguments)

	err := flags.Parse(groups.flags)
	if err != nil {
		return args{}, fmt.Errorf("parse flags: %w", err)
	}

	options.Function = srcloc.FunctionName(function)

	return assignPositionals(options, flags, groups.positional)
}

// Separate flags from positionals: flag.Parse stops at the first positional.
func splitArguments(arguments []string) argumentGroups {
	var groups argumentGroups

	for index := 0; index < len(arguments); index++ {
		argument := arguments[index]

		switch {
		case argument == "--":
			groups.positional = append(groups.positional,
				arguments[index+1:]...)

			return groups
		case argument == "-" || !strings.HasPrefix(argument, "-"):
			groups.positional = append(groups.positional, argument)
		default:
			groups.flags = append(groups.flags, argument)

			if hasFlagValue(arguments, index) {
				index++
				groups.flags = append(groups.flags, arguments[index])
			}
		}
	}

	return groups
}

func hasFlagValue(arguments []string, index int) bool {
	if index+1 >= len(arguments) {
		return false
	}

	switch arguments[index] {
	case "--program", "-program", "--function", "-function":
		return true
	default:
		return false
	}
}

func assignPositionals(
	options args, flags *flag.FlagSet, positional []string,
) (args, error) {
	provided := make(map[string]bool)

	flags.Visit(func(option *flag.Flag) { provided[option.Name] = true })

	if !provided["program"] && len(positional) > 0 {
		options.Program = positional[0]
		positional = positional[1:]
	}

	if !provided["function"] && len(positional) > 0 {
		options.Function = srcloc.FunctionName(positional[0])
		positional = positional[1:]
	}

	if len(positional) > 0 {
		return args{}, fmt.Errorf("%w: %v",
			errUnexpectedPositionals, positional)
	}

	if options.Program == "" {
		return args{}, errEmptyProgram
	}

	if options.Function == "" {
		return args{}, errEmptyFunction
	}

	return options, nil
}
