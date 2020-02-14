package main

import (
	"fmt"
	"os"
)

// OptionOperation command line option
var OptionOperation = Add

// OptionClear command line option
var OptionClear = false

// OptionRelative command line option
var OptionRelative = false

// OptionReadStdin command line option
var OptionReadStdin = false

// OptionWriteStdout command line option
var OptionWriteStdout = false

// OptionVerbose command line option
var OptionVerbose = false

// ArgumentItems command line option
var ArgumentItems = []string{}

// ArgumentExecs command line option
var ArgumentExecs = []string{}

// A list of predefined name and corresponding exec shortcuts
var predefinedExecs = map[string]string{
	"-cp": "cp { } .",
	"-mv": "mv { } .",
	"-rm": "rm { }",
}

// ParseArguments parses the arguments
func ParseArguments() bool {
	args, argsCount := os.Args, len(os.Args)

	allowFiles := true
	onlyFiles := false

	for i := 1; i < argsCount; i++ {
		arg := args[i]
		hasNext := i+1 < argsCount

		if arg[0] != '-' || onlyFiles {
			if !allowFiles {
				return showError("cannot add item '%s' after -exec commands", arg)
			}
			ArgumentItems = append(ArgumentItems, arg)

		} else if arg == "-p" || arg == "-op" || arg == "--op" {
			if hasNext {
				if op, ok := parseOperation(args[i+1]); ok {
					OptionOperation = op
					i++ // TODO export
				} else {
					return showError("parameter '%s' of argument '%s' cannot be parsed", args[i+1], arg)
				}
			} else {
				return showError("argument '%s' needs an parameter", arg)
			}

		} else if arg == "-s" || arg == "-sub" {
			OptionOperation = Remove

		} else if arg == "-x" || arg == "-xor" {
			OptionOperation = Xor

		} else if arg == "-e" || arg == "-exec" || arg == "--exec" {
			if hasNext {
				ArgumentExecs = append(ArgumentExecs, args[i+1])
				i++ // TODO export
				allowFiles = false
			} else {
				return showError("argument '%s' needs an parameter", arg)
			}

		} else if arg == "-" || arg == "-c" || arg == "-clear" || arg == "--clear" {
			OptionClear = true

		} else if arg == "-r" || arg == "-rel" || arg == "--rel" {
			OptionRelative = true

		} else if arg == "-i" || arg == "-in" || arg == "--in" {
			OptionReadStdin = true

		} else if arg == "-o" || arg == "-out" || arg == "--out" {
			OptionWriteStdout = true

		} else if arg == "-v" || arg == "-verbose" || arg == "--verbose" { // TODO add to usage
			OptionVerbose = true

		} else if arg == "-h" || arg == "-help" || arg == "--help" { // TODO add to usage
			return showUsage()

		} else if arg == "--" {
			onlyFiles = true

		} else if exec, ok := predefinedExecs[arg]; ok {
			ArgumentExecs = append(ArgumentExecs, exec)

		} else {
			return showError("unknown argument '%s'", arg)
		}

	}

	return true
}

func parseOperation(text string) (Operation, bool) {
	c := text[0]

	if c == '+' || c == 'a' {
		return Add, true
	} else if c == '-' || c == 's' || c == 'r' || c == 'e' {
		return Remove, true
	} else if c == '^' || c == 'x' || c == 't' {
		return Xor, true
	}
	return 0, false
}

func showError(format string, a ...interface{}) bool {
	fmt.Print("mark: ")
	fmt.Printf(format, a...)
	fmt.Println()
	return false
}

func showUsage() bool {
	fmt.Print(usage)
	return false
}

const usage = `
usage: mark [-] [-op OPERATION] [-in] [-out] item...
       mark [-exec SHELL_COMMAND]...
       mark

 -, -c, -clear, --clear
    Remove all marked items.

 -p, -op, --op OPERATION
    Set what the input items will interact with existing marked items.

    OPERATION

    +, add  (default)  Add the new items to the existing items.
    -, sub, rem, exl   Remove the new items from the existing items.
    ^, xor, tog        Remove its item that exists and add those that don't.

 -s, -sub, -x, -xor
    Shortcuts for '-op sub' and '-op xor'.

 -e, -exec, --exec SHELL_COMMAND
    Execute the shell command with every marked item.

    SHELL_COMMAND

    In the given shell command any '{}' is replaced with an item and any '{ }' is replace with all the items separated by spaces.
    If there is no '{}' or '{ }' and the shell command does not end with a ';', a ' {}' is implied at the end of the provided shell command.
    If only the shell command contains one or more '{ }' and no '{}', it will be executed a single time.

 -cp, -mv, -rm
    Execute the predefined shell commands 'cp { } .', 'mv { } .' and 'rm { }'.

 -i, -in, --in
    Read items form stdin (implicit if a pipe is detected).

 -o, -out, --out
    Output all the marked items to the stdout (implicit if no arguments or pipe is detected).

 -r, -rel, --rel
    The items will not expanded to the absolute path (must be used when the items are not filenames).

`
